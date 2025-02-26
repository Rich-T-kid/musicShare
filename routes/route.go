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
	"time"

	"github.com/Rich-T-kid/musicShare/pkg/logs"
	"github.com/Rich-T-kid/musicShare/pkg/models"
	sw "github.com/Rich-T-kid/musicShare/spotwrapper"
)

var (
	clientID      = "8b277fb167214401bb9486e53d183963"
	clientSecrete = "db97671791ec461f922c52359d89cddf"
	authURL       = "https://accounts.spotify.com/authorize"
	redirect      = "http://18.222.251.123:80/callback"
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

func Callback(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request URL
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Failed to parse request URL", http.StatusInternalServerError)
		log.Printf("Error parsing request URL: %v", err)
		return
	}

	// Extract query parameters
	QV := u.Query()
	code := QV.Get("code")
	state := QV.Get("state")
	spotifyError := QV.Get("error")

	// Validate state to prevent CSRF attacks
	if state != randomString {
		http.Error(w, "Invalid state code returned by Spotify", http.StatusUnauthorized)
		log.Printf("State mismatch: expected %s, got %s", randomString, state)
		return
	}

	// Handle Spotify errors
	if spotifyError != "" {
		http.Error(w, "Spotify login error: "+spotifyError, http.StatusBadRequest)
		log.Printf("Spotify login error: %s", spotifyError)
		return
	}

	// Ensure authorization code is present
	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		log.Println("Authorization code missing in callback")
		return
	}

	// Retrieve access tokens from Spotify
	tokens, err := getToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		log.Printf("Error getting access token: %v", err)
		return
	}

	ctx := context.Background()
	cache := userNamecache

	// Fetch user profile from Spotify
	userProfileData := sw.GetUserData(ctx, tokens.AccessToken)
	if userProfileData == nil {
		http.Error(w, "Failed to fetch user profile data", http.StatusInternalServerError)
		log.Println("Error: user profile data is nil")
		return
	}

	username := userProfileData.DisplayName
	key := fmt.Sprintf("UserName:%s", username)

	// Cache the username if not already present
	if !cache.Exist(ctx, key) {
		cache.Set(ctx, key, "exist", Month*12)
	}
	cache.StoreTokens(username, tokens.AccessToken, tokens.Refresh)

	// Prepare context with username
	ctx = context.WithValue(ctx, models.UsernameKey{}, username)
	userUUIDKey := fmt.Sprintf("UserName:UUID%s", username)
	db := sw.NewDocumentStore()
	fmt.Println("Response to database being connected -> ", db.Connected(ctx))
	// Check if the user exists in cache, otherwise create a new profile
	if !cache.Exist(ctx, userUUIDKey) {
		log.Printf("User %s does not exist, generating new profile", username)

		userDoc, err := sw.NewUserProfile(ctx, tokens.AccessToken)
		if err != nil {
			http.Error(w, "Internal server error while creating user profile", http.StatusInternalServerError)
			log.Printf("Error constructing user MongoDB document for %s: %v", username, err)
			return
		}

		err = sw.SaveUser(userDoc)
		if err != nil {
			http.Error(w, "Internal server error while saving user profile", http.StatusInternalServerError)
			log.Printf("Error saving user MongoDB document for %s: %v", username, err)
			return
		}

		cache.Set(ctx, userUUIDKey, userDoc.UUID, Month*12)
		log.Printf("User %s's MongoDB document was generated with UUID: %s", userDoc.UserProfileResponse.DisplayName, userDoc.UUID)
	}

	// Retrieve the user's UUID from the cache
	log.Println("Attempting to retrieve user from Redis cache")
	presentUserUUID := cache.Get(ctx, userUUIDKey)

	if presentUserUUID == "" {
		http.Error(w, "User UUID not found in cache", http.StatusInternalServerError)
		log.Printf("Error: User UUID missing from cache for %s", username)
		return
	}

	log.Printf("User UUID retrieved from cache: %s", presentUserUUID)

	// Fetch user document from MongoDB
	userDoc, err := sw.GetUserByID(presentUserUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("User %s not found in MongoDB", username), http.StatusNotFound)
		log.Printf("Error retrieving user %s with UUID %s: %v", username, presentUserUUID, err)
		return
	}
	if userDoc == nil {
		http.Error(w, "User document not found", http.StatusNotFound)
		log.Printf("Error: Retrieved user document is nil for UUID: %s", presentUserUUID)
		return
	}

	// Prepare and send response
	userinfo := struct {
		UUID             string
		Name             string
		UserCreationDate time.Time
		Date             time.Time
	}{
		UUID:             userDoc.UUID,
		Name:             userDoc.UserProfileResponse.DisplayName,
		UserCreationDate: userDoc.CreatedAt,
		Date:             time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func LoveShare(w http.ResponseWriter, r *http.Request) {

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
	w.Header().Set("Content-Type", "application/json")
	var baseURL = "https://accounts.spotify.com/authorize"
	var username = randomString
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
func getToken(code string) (models.TokenResponse, error) {
	var endpoint = "https://accounts.spotify.com/api/token"
	var invalidResponse models.TokenResponse

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
		fmt.Println("Error creating request:", err)
		return invalidResponse, err
	}

	// Set correct headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Printf("Exact request sent to spotify %+v", resp.Request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return invalidResponse, err
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Token Response Status: %d\n,  Token Response Body %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return invalidResponse, fmt.Errorf("error: response status code %d", resp.StatusCode)
	}

	var response models.TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return invalidResponse, fmt.Errorf("error parsing response: %v", err)
	}

	return response, nil
}
