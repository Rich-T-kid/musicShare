package routes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"

	"loveShare/logs"
	sw "loveShare/spotWrapper"
)

var (
	clientID      = "8b277fb167214401bb9486e53d183963"
	clientSecrete = "db97671791ec461f922c52359d89cddf"
	authURL       = "https://accounts.spotify.com/authorize"
	redirect      = "http://localhost:8080/callback"
	token_url     = "https://accounts.spotify.com/api/token"
	scopes        = "user-library-read user-modify-playback-state playlist-modify-public playlist-modify-private playlist-read-private user-top-read user-follow-read"
	randomString  = "ChangeLater"
	IDComboHash64 = "OGIyNzdmYjE2NzIxNDQwMWJiOTQ4NmU1M2QxODM5NjM6ZGI5NzY3MTc5MWVjNDYxZjkyMmM1MjM1OWQ4OWNkZGY="
	userNamecache = sw.NewCache[string, string]()
	logger        = logs.NewLogger() //TODO:
	Month         = 720
)

// HTML template rendering
func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		logger.Warning(fmt.Sprintf("Error reading template:  %e", err))
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error Executing template:  %e", err))
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}

}

func RedirectPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Info(fmt.Sprintf("Error reading template:  %e", err))
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("Error Executing template:  %e", err))
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}

}

// html template
func SongofDay(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	QV := u.Query()
	code := QV.Get("code")
	username := QV.Get("state")
	spotifyError := QV.Get("error")
	if spotifyError != "" {
		http.Error(w, "Error occurred with Spotify login", http.StatusBadRequest)
		logger.Critical(fmt.Sprintf("Spotify Returned an error instead of a valid code %e", err))
		return
	}

	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	tokens, err := getToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	cache := userNamecache
	key := fmt.Sprintf("UniqueUserName:%s", username)
	if !cache.Exist(ctx, key) {
		cache.Set(ctx, key, "1", Month)
	}

	ctx = context.WithValue(ctx, sw.UsernameKey{}, username) // Username is passed along to all request made here
	cache.StoreTokens(username, tokens.AccessToken, tokens.Refresh)
	res, err := sw.NewUserProfile(ctx, tokens.AccessToken)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info(fmt.Sprintf("Profile Response %v", res))
	encoder := json.NewEncoder(w)
	encoder.Encode(res)

}

func LoveShare(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/SongofDay.html")
	if err != nil {
		logger.Info(fmt.Sprintf("Error reading template:  %e", err))
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("Error executing template:  %e", err))
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// returns link that user will use to login with spotify
func RedirectLink(w http.ResponseWriter, r *http.Request) {
	var baseURL = "https://accounts.spotify.com/authorize"
	var username = r.Header.Get("X-username")
	if username == "" {
		w.WriteHeader(400)
		w.Write([]byte("Username must be provider in the headers as X-username : {Username}"))
		return
	}
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirect)
	// in the future the "state" that we pass along will be the username that is provided when calling this endpoint
	params.Add("state", username)
	params.Add("scope", scopes)
	params.Add("show_dialog", "true")

	// Append query parameters to the URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"link": fullURL,
	}
	w.WriteHeader(200)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // Prevent & from becoming \u0026
	encoder.Encode(response)
}

// handles adding user to databse as well as tracking when they logged in.
// if succseful returns 200 status code and user can procced to the -> login with spofify
func SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		sw.SaveUser()
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
	default:
		w.Write([]byte("Wrong method type, Must be Post Request"))
		w.WriteHeader(405)
	}
}

// home page -> http request to /signIn if 200 response -> home page/login with spotify, if 200 response -> redirect page where song of the day is

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}
func getToken(code string) (sw.TokenResponse, error) {
	var endpoint = "https://accounts.spotify.com/api/token"
	var invalidResponse sw.TokenResponse

	// Correctly format client credentials
	idSecretCombo := fmt.Sprintf("%s:%s", clientID, clientSecrete)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(idSecretCombo))

	// Properly encode request body for application/x-www-form-urlencoded
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirect)
	encodedBody := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", endpoint, encodedBody)
	if err != nil {
		logger.Route(fmt.Sprintf("Error creating request %e", err))
		return invalidResponse, err
	}

	// Set correct headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Route(fmt.Sprintf("Error making request: %e", err))
		return invalidResponse, err
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)
	logger.Info(fmt.Sprintf("Response Status: %d\nResponse Body: %s\n", resp.StatusCode, string(body)))

	if resp.StatusCode != http.StatusOK {
		logger.Warning(fmt.Sprintf("Generating Token response for code %s resulted in a non 200 status code %d", code, resp.StatusCode))
		return invalidResponse, fmt.Errorf("error: response status code %d", resp.StatusCode)
	}

	var response sw.TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return invalidResponse, fmt.Errorf("error parsing response: %v", err)
	}

	return response, nil
}

// just returns a access token and that tokens expiration time

type Response struct {
	Message string `json:"message"`
}

func UniqueUsername(w http.ResponseWriter, r *http.Request) {
	cache := userNamecache
	w.Header().Set("Content-Type", "application/json") // Ensure JSON response

	vars := mux.Vars(r)
	name, exists := vars["name"]
	if !exists {
		http.Error(w, `{"message": "Missing parameter"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	validUsername := cache.Get(ctx, fmt.Sprintf("UniqueUserName:%s", name))

	var resp Response
	if validUsername == "" {
		w.WriteHeader(http.StatusOK)
		resp.Message = fmt.Sprintf("Username %s is available", name)
	} else {
		w.WriteHeader(http.StatusConflict)
		logger.Info(fmt.Sprintf("Username %s is already taken, must select new one", name))
		resp.Message = fmt.Sprintf("Username %s already exists. Choose another one", name)
	}

	json.NewEncoder(w).Encode(resp) // Encode response as JSON
}
