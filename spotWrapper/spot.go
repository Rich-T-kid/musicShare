package spotwrapper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// handles the API related stuff

// Base Response Interface

// FetchSpotifyTop fetches top artists or tracks from Spotify API
func FetchSpotifyTop(ctx context.Context, accessToken string, dataType string) (SpotifyTopResponse, error) {
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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
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

	case http.StatusBadRequest: // 400
		return nil, errors.New("bad request: check your request parameters")
	case http.StatusUnauthorized: // 401
		return nil, errors.New("unauthorized: invalid or expired access token")
	case http.StatusForbidden: // 403
		return nil, errors.New("forbidden: you do not have permission to access this resource")
	case http.StatusNotFound: // 404
		return nil, errors.New("not found: the requested resource does not exist")
	case http.StatusTooManyRequests: // 429
		return nil, errors.New("rate limited: too many requests, try again later")
		// this needs to be handled explicitly. im thinking pass a overload struct to the ctx and handle it here as well
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
