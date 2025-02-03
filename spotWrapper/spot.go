package spotwrapper

import (
	"fmt"
	"io"
	"net/http"
)

// handles the API related stuff

func FetchSpotifyTop(accessToken string, dataType string) (string, error) {
	fmt.Println("access token", accessToken)
	// Validate input type
	if dataType != "artists" && dataType != "tracks" {
		return "", fmt.Errorf("invalid type: must be 'artists' or 'tracks'")
	}

	// Construct the Spotify API URL
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/%s", dataType)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Return the response as a string
	return string(body), nil
}
