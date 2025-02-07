package spotwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

// Helper function to handle errors consistently
func handleError(err error, context string) {
	if err != nil {
		log.Fatalf("Error in %s: %v", context, err)
	}
}

// Helper function to check if a response is successful
func checkResponseStatusCode(resp *http.Response, validCodes []int) error {
	for _, code := range validCodes {
		if resp.StatusCode == code {
			return nil
		}
	}
	return fmt.Errorf("unexpected status code %d", resp.StatusCode)
}

// Have the overloader perform all the request after you test each method
// Function to get user data and unmarshal it into the provided struct
// Works
func GetUserData(token string) *UserProfileResponse {
	// Hardcode the endpoint for testing purposes
	endpoint := "https://api.spotify.com/v1/me"

	// Validate inputs
	if token == "" {
		handleError(fmt.Errorf("access token is empty"), "getUserData")
	}

	// Create the request
	req, err := http.NewRequest("GET", endpoint, nil)
	handleError(err, "http.NewRequest")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	handleError(err, "http.DefaultClient.Do")
	defer resp.Body.Close()

	// Check for expired access token
	if resp.StatusCode == 401 {
		log.Fatal("Access Token has expired. You need to grab a new one.")
	}

	// Validate the response status code
	err = checkResponseStatusCode(resp, []int{200})
	handleError(err, "checkResponseStatusCode")
	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	handleError(err, "io.ReadAll")
	// Unmarshal the response body into the provided struct
	var dest UserProfileResponse
	err = json.Unmarshal(bodyBytes, &dest)
	handleError(err, "json.Unmarshal")
	return &dest
}

// ConvertToFollowedArtists converts SpotArtist struct to FollowedArtist structs
func ConvertToFollowedArtists(spotArtists *SpotArtist) []FollowedArtist {
	if spotArtists == nil {
		handleError(fmt.Errorf("spotArtists is nil"), "ConvertToFollowedArtists")
	}

	var followedArtists []FollowedArtist

	// Iterate over the items in SpotArtist
	for _, artist := range spotArtists.Artists.Items {
		// Ensure valid data before using it
		if artist.Name == "" || artist.ExternalUrls.Spotify == "" || artist.URI == "" {
			log.Printf("Skipping invalid artist: %+v", artist)
			continue
		}

		// Create FollowedArtist from SpotArtist fields
		followedArtist := FollowedArtist{
			Name:    artist.Name,
			Spotify: artist.ExternalUrls.Spotify,
			Genres:  artist.Genres,
			URI:     artist.URI,
		}

		// Append to the result
		followedArtists = append(followedArtists, followedArtist)
	}

	return followedArtists
}

// Function to get artist information and convert it to FollowedArtist structs
func ArtistInfo(token string) []FollowedArtist {
	// Hardcode the endpoint for testing purposes
	endpoint := "https://api.spotify.com/v1/me/following?type=artist"

	// Validate inputs
	if token == "" {
		handleError(fmt.Errorf("access token is empty"), "ArtistInfo")
	}

	// Create the request
	req, err := http.NewRequest("GET", endpoint, nil)
	handleError(err, "http.NewRequest")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	handleError(err, "http.DefaultClient.Do")
	defer resp.Body.Close()

	// Check for expired access token
	if resp.StatusCode == 401 {
		log.Fatal("Access Token has expired. You need to grab a new one.")
	}

	// Validate the response status code
	err = checkResponseStatusCode(resp, []int{200})
	handleError(err, "checkResponseStatusCode")

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	handleError(err, "io.ReadAll")
	var dest SpotArtist
	// Unmarshal the response body into the provided struct
	err = json.Unmarshal(bodyBytes, &dest)
	handleError(err, "json.Unmarshal")

	// Return the converted FollowedArtist structs
	return ConvertToFollowedArtists(&dest)
}

// Function to create a playlist using the Spotify API
func CreatePlaylist(token, spotifyID, playlistName, description string) (PlaylistResponse, error) {
	// Hardcode the endpoint for testing purposes
	endpoint := "https://api.spotify.com/v1/users/" + spotifyID + "/playlists"

	// Validate inputs
	if token == "" || spotifyID == "" || playlistName == "" {
		handleError(fmt.Errorf("missing required parameter(s)"), "CreatePlaylist")
	}

	// Prepare request body
	body := CreatePlaylistRequestBody{
		Name:        playlistName,
		Description: description,
		Public:      false,
	}

	// Marshal the struct into JSON format
	jsonData, err := json.Marshal(body)
	handleError(err, "json.Marshal")

	// Create the request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	handleError(err, "http.NewRequest")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	handleError(err, "http.DefaultClient.Do")
	defer resp.Body.Close()

	// Validate the response status code
	err = checkResponseStatusCode(resp, []int{201})
	if err != nil {
		// If not 201, log the response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response body: %s", string(bodyBytes))
		return PlaylistResponse{}, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	handleError(err, "io.ReadAll")

	// Unmarshal the response into PlaylistResponse struct
	var response PlaylistResponse
	err = json.Unmarshal(bodyBytes, &response)
	handleError(err, "json.Unmarshal")

	// Return the response
	return response, nil
}

// Function to add a track to a playlist
func AddToPlaylist(token, songURI, playlistID string) bool {
	// Hardcode the endpoint for testing purposes
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	// Validate inputs
	if token == "" || songURI == "" || playlistID == "" {
		handleError(fmt.Errorf("missing required parameter(s)"), "AddToPlaylist")
	}

	// Prepare the request body
	body := PlaylistAddtionRequest{
		URI:      []string{songURI},
		Position: 0,
	}

	// Marshal the struct into JSON format
	jsonData, err := json.Marshal(body)
	handleError(err, "json.Marshal")

	// Create the request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	handleError(err, "http.NewRequest")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	handleError(err, "http.DefaultClient.Do")
	defer resp.Body.Close()

	// Validate the response status code
	err = checkResponseStatusCode(resp, []int{200, 201})
	if err != nil {
		// Log any error response
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response body: %s", string(bodyBytes))
		log.Fatal(err)
	}

	// Return true if successful
	return true
}

// slight variation to the above method. Only difference is that the song URI doesnt need to have the spotify:track: prefix
// more than likely will remove
func addtoPlaylist(endpoint, token, songURI, playlistID string) bool {
	// Ensure valid inputs
	if endpoint == "" || token == "" || songURI == "" || playlistID == "" {
		log.Println("Error: Missing required parameters (endpoint, token, songURI, playlistID)")
		return false
	}

	// Format the endpoint for adding tracks to a playlist
	endpoint = fmt.Sprintf("%s%s/tracks", endpoint, playlistID)

	// Prepare the request body
	body := PlaylistAddtionRequest{
		URI:      []string{fmt.Sprintf("spotify:track:%s", songURI)},
		Position: 0,
	}

	// Marshal the body into JSON format
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return false
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return false
	}

	// Set headers for the request
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing HTTP request: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Handle response status codes
	if resp.StatusCode == 401 {
		log.Println("Error: Access Token has expired. Please grab a new one.")
		return false
	}

	if resp.StatusCode == 201 {
		log.Println("Track added to playlist successfully.")
		return true
	}

	// Handle other unexpected status codes
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return false
	}

	// Log response body and status code for debugging purposes
	log.Printf("Unexpected status code: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(bodyBytes))

	return false
}

// Function to parse artist data and return a list of UserTopArtist structs
func parseArtist(data SpotifyTopArtistResponse) []UserTopArtist {
	var result []UserTopArtist

	for _, item := range data.Items {
		// Check if artist image exists or not
		artistPhoto := "N/A"
		if len(item.Images) > 0 {
			artistPhoto = item.Images[0].URL
		}

		// Create UserTopArtist and append to result
		userArtist := UserTopArtist{
			Name:        item.Name,
			URI:         item.URI,
			Genres:      item.Genres,
			ArtistPhoto: artistPhoto,
		}
		result = append(result, userArtist)
	}
	return result
}

// Function to fetch and process the user's top artists
func TopArtist(token string) []UserTopArtist {
	endpoint := "https://api.spotify.com/v1/me/top/artists"
	if token == "" {
		handleError(fmt.Errorf("access token is empty"), "TopArtist")
	}

	// Initialize variables
	var allArtists []UserTopArtist
	limit := 5 // Set a limit on the number of requests to avoid infinite loops
	pageCount := 0

	for {
		// Check if we've reached the limit for the number of requests
		if pageCount >= limit {
			log.Printf("Reached the request limit of %d pages, stopping further requests.\n", limit)
			break
		}

		// Create the request
		req, err := http.NewRequest("GET", endpoint, nil)
		handleError(err, "http.NewRequest")
		req.Header.Set("Authorization", "Bearer "+token)

		// Execute the request
		resp, err := http.DefaultClient.Do(req)
		handleError(err, "http.DefaultClient.Do")
		defer resp.Body.Close()

		// Check for expired access token
		if resp.StatusCode == 401 {
			log.Fatal("Access Token has expired. You need to grab a new one.")
		}

		// Validate the response status code
		err = checkResponseStatusCode(resp, []int{200})
		handleError(err, "checkResponseStatusCode")

		// Read the response body
		bodyBytes, err := io.ReadAll(resp.Body)
		handleError(err, "io.ReadAll")

		// Unmarshal the response body into the SpotifyTopArtistResponse struct
		var dest SpotifyTopArtistResponse
		err = json.Unmarshal(bodyBytes, &dest)
		handleError(err, "json.Unmarshal")

		// Parse the artist data and append it to the result
		artists := parseArtist(dest)
		allArtists = append(allArtists, artists...)

		// If there's a next page, update the endpoint
		if dest.Next == "" {
			break // No more data to fetch, exit the loop
		}

		// Update the endpoint to the next page
		endpoint = dest.Next

		// Increment the request count
		pageCount++
	}

	// Return the full list of artists
	return allArtists
}

func parseTracks(response *SpotifyTrackResponse, topTrack *UserTopTrack) {
	// Parse top albums
	for _, item := range response.Items {
		album := Album{
			Artist:      item.Artists[0].URI,
			Name:        item.Album.Name,
			AlbumLink:   fmt.Sprintf("https://spotify.com/album/%s", item.Album.URI),
			AlbumURI:    item.Album.URI,
			AlbumID:     item.ID,
			AlbumImage:  Image{URL: item.Album.Images[0].URL},
			AlbumName:   item.Album.Name,
			TotalTracks: item.Album.TotalTracks,
			ReleaseDate: item.Album.ReleaseDate,
		}
		// Append to existing slice
		topTrack.TopAlbums = append(topTrack.TopAlbums, album)

		// Parse top singles (tracks)
		track := SingleTrack{
			Artist:      item.Artists[0].URI,
			Name:        item.Name,
			TrackLink:   fmt.Sprintf("https://spotify.com/track/%s", item.URI),
			TrackName:   item.Name,
			ReleaseDate: item.Album.ReleaseDate, // Assuming same release date for simplicity
		}
		// Append to existing slice
		topTrack.TopSingles = append(topTrack.TopSingles, track)
	}
}

func TopTracks(token string) UserTopTrack {
	endpoint := "https://api.spotify.com/v1/me/top/tracks"
	var result UserTopTrack
	limit := 50 // Set limit to 50 to fetch 50 items per page (Spotify API max)
	currentPage := 0
	maxpages := 3

	for currentPage <= maxpages {
		// Create the request
		req, err := http.NewRequest("GET", endpoint, nil)
		handleError(err, "http.NewRequest")
		req.Header.Set("Authorization", "Bearer "+token)
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", limit)) // Ensure the limit is set to 50
		req.URL.RawQuery = q.Encode()

		// Execute the request
		resp, err := http.DefaultClient.Do(req)
		handleError(err, "http.DefaultClient.Do")
		defer resp.Body.Close()

		// Check for expired access token
		if resp.StatusCode == 401 {
			log.Fatal("Access Token has expired. You need to grab a new one.")
		}

		// Validate the response status code
		err = checkResponseStatusCode(resp, []int{200})
		handleError(err, "checkResponseStatusCode")

		// Read the response body
		bodyBytes, err := io.ReadAll(resp.Body)
		handleError(err, "io.ReadAll")

		// Unmarshal the response body into the SpotifyTrackResponse struct
		var response SpotifyTrackResponse
		err = json.Unmarshal(bodyBytes, &response)
		handleError(err, "json.Unmarshal")

		// Pass pointer to the result to append items
		parseTracks(&response, &result)

		// If there is no next page, break the loop
		if response.Next == "" || currentPage >= maxpages {
			break
		}

		// Update the endpoint to the next page for pagination
		endpoint = response.Next
		currentPage++
	}

	// Debug: Print the length of TopAlbums and TopSingles

	return result
}

func NewDocument(followed []FollowedArtist, topTracks UserTopTrack, UserFavorites []UserTopArtist) *UserMusicInfo {
	return &UserMusicInfo{
		FollowedArtist: followed,
		TopTracks:      topTracks,
		TopsArtist:     UserFavorites,
	}
}

func NewUserProfile(token, userID string) *UserMusicInfo {
	currentTime := time.Now()
	fmt.Printf("Finished proccessing new user %s at %v", userID, currentTime)
	return NewDocument(ArtistInfo(token), TopTracks(token), TopArtist(token))
}
