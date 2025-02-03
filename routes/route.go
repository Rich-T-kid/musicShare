package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

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
)

// HTML template rendering
func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println("Error reading template ->>", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error reading template ->>", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}

}

func RedirectPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Error reading template ->>", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error reading template ->>", err)
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
	//state := QV.Get("state") this should just be a username for now so we dont need it
	spotifyError := QV.Get("error")
	if spotifyError != "" {
		http.Error(w, "Error occurred with Spotify login", http.StatusBadRequest)
		fmt.Println("Spotify callback error:", spotifyError)
		return
	}

	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}
	fmt.Println("Code recived ->", code)
	// Step 1: Exchange code for tokens
	//TODO: This is the most imprtant for now -> this needs to sit behind a cache
	/*
		tokens, err := getToken(code)
		if err != nil {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			fmt.Println("Error getting access token:", err)
			return
		}*/

	// Step 2: Store tokens in MongoDB (Pseudocode)
	// TODO: Implement actual MongoDB storage
	//fmt.Println("Storing tokens in DB: AccessToken:", tokens.AccessToken, "RefreshToken:", tokens.Refresh)
	res, _ := sw.FetchSpotifyTop("BQDV9UxmZ1llhnBiJEQ8QAJG_A4sG-tOQDrBVRndZbquzVLyn6nq6dw5sGmPfoLSQupQTIE9UfMtZ6xBzoLPkYaxrfHuaXT0-BUMEO9CTgNO4glGpEx3HDYcJAgw1CSrs6JXoBadJdBuwgh1hJxnortvJSSwDtQKvk8Zowe43bZy_xKbDlYPkbp8JEgfakQWwgZ9_7Nlp5rvYc5YXxYimGzu40UfH4zmkp67kMhU3DhH5Dgl1zFy1jy-YE9HZOSQLCHJlw-OxM_SF8YYCfAPvt5RkFLme9Q3OnuPuqUmSzbVekhGxRX0-Unz7uRKvhy3VTAP", "tracks")
	fmt.Println("Result ,", res)
	//TODO: this is where id like store the tokens in mongo DB so we dont need to always look it up
	tmpl, err := template.ParseFiles("templates/SongofDay.html")
	if err != nil {
		fmt.Println("Error reading template ->>", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template without any variables (blank for now)
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error reading template ->>", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// returns link that user will use to login with spotify
func RedirectLink(w http.ResponseWriter, r *http.Request) {
	var baseURL = "https://accounts.spotify.com/authorize"
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirect)
	// in the future the "state" that we pass along will be the username that is provided when calling this endpoint
	params.Add("state", randomString)
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

// workig with song of the day stuff. least important right now
func Song(w http.ResponseWriter, r *http.Request) {}

// home page -> http request to /signIn if 200 response -> home page/login with spotify, if 200 response -> redirect page where song of the day is

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}

func getToken(code string) (tokenResponse, error) {
	var endpoint = "https://accounts.spotify.com/api/token"
	var invalidResponse tokenResponse
	//idSecreteCombo := fmt.Sprintf("%s:%s", clientID, clientSecrete)
	headerAuth := IDComboHash64 //base64.StdEncoding.EncodeToString([]byte(idSecreteCombo))

	body := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": redirect,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return invalidResponse, fmt.Errorf("error marshalling spotify info into json %e", err)
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return invalidResponse, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", headerAuth)) // Example token header

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return invalidResponse, err
	}
	defer resp.Body.Close()
	var Response tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&Response)
	if err != nil {
		return invalidResponse, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("ERROR: response status code %d , response body %v\n", resp.StatusCode, resp.Body)
	}
	fmt.Printf("status code : %d json Response from token endpoint %v\n", resp.StatusCode, Response)
	return Response, nil
}

// just returns a access token and that tokens expiration time
func refresh(refreshToken string) (string, int, error) {
	var endpoint = "https://accounts.spotify.com/api/token"

	// Encode clientID:clientSecret in Base64
	//idSecretCombo := fmt.Sprintf("%s:%s", clientID, clientSecrete)
	headerAuth := IDComboHash64 //base64.StdEncoding.EncodeToString([]byte(idSecretCombo))

	// Form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+headerAuth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Decode response
	var Response refreshResponse
	err = json.NewDecoder(resp.Body).Decode(&Response)
	if err != nil {
		return "", 0, fmt.Errorf("error decoding response: %w", err)
	}

	return Response.AccessToken, Response.Expir, nil
}
