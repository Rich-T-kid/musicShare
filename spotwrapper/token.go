package spotwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Rich-T-kid/musicShare/pkg"
)

// handle everything associated with tokens. Should be pretty small and low effort work since its already been done in route.go
// focus on utilty functions as well as storing tokens in a datastore(reddis for now, mongoDB later), as well as utility function/endpoint that will showcase all users access and refresh token for testing
var (
	tokenStore    = NewCache[string, string]()
	IDComboHash64 = "OGIyNzdmYjE2NzIxNDQwMWJiOTQ4NmU1M2QxODM5NjM6ZGI5NzY3MTc5MWVjNDYxZjkyMmM1MjM1OWQ4OWNkZGY="
)

// handleing 401 responses
func refreshEndPoint(refreshToken string) (*pkg.RefreshResponse, error) {
	endpoint := "https://accounts.spotify.com/api/token"

	headerAuth := IDComboHash64 // Assumes `IDComboHash64` is correctly set

	// Form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+headerAuth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Decode response
	var response pkg.RefreshResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// ✅ Handle API errors like invalid refresh tokens
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("spotify api returned status %d: %v", resp.StatusCode, response)
	}

	return &response, nil
}

func validateAccessToken(accessToken string) bool {
	// for now were just requesting the spotify api. in the future we should be able to just check the cache first
	var baseURL = "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", baseURL, nil)

	if err != nil {
		return false
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode != 401 && resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func generateAccessToken(refreshT string) string {
	RefreshResponse, err := refreshEndPoint(refreshT)
	if err != nil {
		log.Fatal(err)
	}
	return RefreshResponse.AccessToken
}

func contains(matches []string, word string) bool {
	for _, words := range matches {
		if words == word {
			return true
		}
	}
	return false
}
