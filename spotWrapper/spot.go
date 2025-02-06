package spotwrapper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// handles the API related stuff

// Base Response Interface

var (
	proxy = newOverloader()
)

// FetchSpotifyTop fetches top artists or tracks from Spotify API
func FetchSpotifyTop(ctx context.Context, userid, accessToken, dataType string) (SpotifyTopResponse, error) {
	// Validate input type
	if dataType != "artists" && dataType != "tracks" {
		return nil, errors.New("invalid type: must be 'artists' or 'tracks'")
	}
	// Construct the Spotify API URL
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/%s", dataType)

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request with timeout
	resp, err := proxy.RetryRequest(ctx, req, userid)
	if err != nil {
		fmt.Println("1")
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Handle different response statuses
	switch resp.StatusCode {
	case http.StatusOK: // 200
		// Decode JSON into the appropriate struct based on dataType
		if dataType == "artists" {
			var response UserTopArtist
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return nil, fmt.Errorf("error decoding JSON: %w", err)
			}
			return &response, nil
		} else {
			var response UserResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return nil, fmt.Errorf("error decoding JSON: %w", err)
			}
			return &response, nil
		}

	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
