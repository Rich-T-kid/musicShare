package spotwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Rich-T-kid/musicShare/pkg"
	"github.com/Rich-T-kid/musicShare/pkg/models"
)

// handles the API related stuff

// Base Response Interface

var (
	proxy = newOverloader()
)

const PlayListName = "Rhythm Reflections"
const PlaylistDescription = "This playlist is your personal musical corner, capturing every track that resonates with you. Each like you give adds a new gem, building a reflection of your evolving taste. Whether it’s a dance anthem or a calming instrumental, watch your collection grow with every fresh discovery. Let these melodies spark joyful memories and inspire new favorites. Enjoy an ever-evolving soundtrack that tells your story, one beloved track at a time!"

// Helper function to handle errors consistently
func handleError(err error, context string) {
	if err != nil {
		log.Fatalf("Error in %s: ->>> %v", context, err)
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
func GetUserData(ctx context.Context, token string) *models.UserProfileResponse {
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
	//  WARNING: DO not replace with proxy. need to grab username before setting it in context
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
	var dest models.UserProfileResponse
	err = json.Unmarshal(bodyBytes, &dest)
	handleError(err, "json.Unmarshal")
	return &dest
}

// ConvertToFollowedArtists converts SpotArtist struct to FollowedArtist structs
func ConvertToFollowedArtists(spotArtists *models.SpotArtist) []models.FollowedArtist {
	if spotArtists == nil {
		handleError(fmt.Errorf("spotArtists is nil"), "ConvertToFollowedArtists")
	}

	var followedArtists []models.FollowedArtist

	// Iterate over the items in SpotArtist
	for _, artist := range spotArtists.Artists.Items {
		// Ensure valid data before using it
		if artist.Name == "" || artist.ExternalUrls.Spotify == "" || artist.URI == "" {
			log.Printf("Skipping invalid artist: %+v", artist)
			continue
		}

		// Create FollowedArtist from SpotArtist fields
		followedArtist := models.FollowedArtist{
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
func ArtistInfo(ctx context.Context, token string) []models.FollowedArtist {
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
	resp, err := proxy.RetryRequest(ctx, req)
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
	var dest models.SpotArtist
	// Unmarshal the response body into the provided struct
	err = json.Unmarshal(bodyBytes, &dest)
	handleError(err, "json.Unmarshal")

	// Return the converted FollowedArtist structs
	return ConvertToFollowedArtists(&dest)
}

// Function to create a playlist using the Spotify API
func CreatePlaylist(ctx context.Context, token, spotifyID, playlistName, description string) (models.PlaylistResponse, error) {
	// Hardcode the endpoint for testing purposes
	endpoint := "https://api.spotify.com/v1/users/" + spotifyID + "/playlists"

	// Validate inputs
	if token == "" || spotifyID == "" || playlistName == "" {
		handleError(fmt.Errorf("missing required parameter(s)"), "CreatePlaylist")
	}

	// Prepare request body
	body := models.CreatePlaylistRequestBody{
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
	resp, err := proxy.RetryRequest(ctx, req)
	handleError(err, "http.DefaultClient.Do")
	defer resp.Body.Close()

	// Validate the response status code
	err = checkResponseStatusCode(resp, []int{201})
	if err != nil {
		// If not 201, log the response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response body: %s", string(bodyBytes))
		return models.PlaylistResponse{}, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	handleError(err, "io.ReadAll")

	// Unmarshal the response into PlaylistResponse struct
	var response models.PlaylistResponse
	err = json.Unmarshal(bodyBytes, &response)
	handleError(err, "json.Unmarshal")

	// Return the response
	return response, nil
}

// Function to add a track to a playlist
func AddToPlaylist(ctx context.Context, token, songURI, playlistID string) bool {
	// Hardcode the endpoint for testing purposes
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	// Validate inputs
	if token == "" || songURI == "" || playlistID == "" {
		handleError(fmt.Errorf("missing required parameter(s)"), "AddToPlaylist")
	}

	// Prepare the request body
	body := models.PlaylistAddtionRequest{
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
	resp, err := proxy.RetryRequest(ctx, req)
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
func addtoPlaylist(ctx context.Context, endpoint, token, songURI, playlistID string) bool {
	// Ensure valid inputs
	if endpoint == "" || token == "" || songURI == "" || playlistID == "" {
		log.Println("Error: Missing required parameters (endpoint, token, songURI, playlistID)")
		return false
	}

	// Format the endpoint for adding tracks to a playlist
	endpoint = fmt.Sprintf("%s%s/tracks", endpoint, playlistID)

	// Prepare the request body
	body := models.PlaylistAddtionRequest{
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
	resp, err := proxy.RetryRequest(ctx, req)
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
func parseArtist(data models.SpotifyTopArtistResponse) []models.UserTopArtist {
	var result []models.UserTopArtist

	for _, item := range data.Items {
		// Check if artist image exists or not
		artistPhoto := "N/A"
		if len(item.Images) > 0 {
			artistPhoto = item.Images[0].URL
		}

		// Create UserTopArtist and append to result
		userArtist := models.UserTopArtist{
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
func TopArtist(ctx context.Context, token string) []models.UserTopArtist {
	endpoint := "https://api.spotify.com/v1/me/top/artists"
	if token == "" {
		handleError(fmt.Errorf("access token is empty"), "TopArtist")
	}

	// Initialize variables
	var allArtists []models.UserTopArtist
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
		resp, err := proxy.RetryRequest(ctx, req)
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
		var dest models.SpotifyTopArtistResponse
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

func parseTracks(response *models.SpotifyTrackResponse, topTrack *models.UserTopTrack) {
	// Parse top albums
	for _, item := range response.Items {
		album := models.Album{
			Artist:      item.Artists[0].URI,
			Name:        item.Album.Name,
			AlbumLink:   fmt.Sprintf("https://spotify.com/album/%s", item.Album.URI),
			AlbumURI:    item.Album.URI,
			AlbumID:     item.ID,
			AlbumImage:  models.Image{URL: item.Album.Images[0].URL},
			AlbumName:   item.Album.Name,
			TotalTracks: item.Album.TotalTracks,
			ReleaseDate: item.Album.ReleaseDate,
		}
		// Append to existing slice
		topTrack.TopAlbums = append(topTrack.TopAlbums, album)

		// Parse top singles (tracks)
		track := models.SingleTrack{
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

func TopTracks(ctx context.Context, token string) models.UserTopTrack {
	endpoint := "https://api.spotify.com/v1/me/top/tracks"
	var result models.UserTopTrack
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
		resp, err := proxy.RetryRequest(ctx, req)
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
		var response models.SpotifyTrackResponse
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

func NewDocument(followed []models.FollowedArtist, topTracks models.UserTopTrack, UserFavorites []models.UserTopArtist) *models.UserMusicInfo {
	return &models.UserMusicInfo{
		FollowedArtist: followed,
		TopTracks:      topTracks,
		TopsArtist:     UserFavorites,
	}
}

// Starts off as empty doesnt need to be allocated right now
func NewMusicPlaylist(playlistID string) *models.MusicSharePlaylist {
	return &models.MusicSharePlaylist{
		Name:        PlayListName,
		PlaylistURI: playlistID,
		Songs:       make([]models.SpotifyURI, 0),
	}
}
func NewDBDocument(userProfile models.UserProfileResponse, userMusicInfo models.UserMusicInfo, playlist models.MusicSharePlaylist) *models.UserMongoDocument {
	return &models.UserMongoDocument{
		UUID:                pkg.NewUUID(),
		UserProfileResponse: userProfile,
		UserMusicInfo:       userMusicInfo,
		MusicSharePlaylist:  playlist,
		Comments:            make([]models.UserComments, 0),
		LikedSongs:          make([]models.SpotifyURI, 0),
		DislikedSongs:       make([]models.SpotifyURI, 0),
		CreatedAt:           time.Now(),
		Updated:             time.Now(),
	}
}

// NewUserProfile builds a new user profile document by fetching music data,
// user data, and creating a new playlist.
func NewUserProfile(ctx context.Context, token string) (*models.UserMongoDocument, error) {
	// Record current time, which can be useful for logging
	currentTime := time.Now()

	// Retrieve the user ID from context
	userID, ok := ctx.Value(models.UsernameKey{}).(string)
	if !ok {
		return nil, fmt.Errorf("username was not properly set in the context")
	}
	fmt.Printf("Finished processing new user %s at %v\n", userID, currentTime)

	// Gather music data from Spotify
	// Note: These functions (ArtistInfo, TopTracks, TopArtist) do not return errors,
	// so we call them directly within NewDocument.
	userMusicInfo := NewDocument(
		ArtistInfo(ctx, token),
		TopTracks(ctx, token),
		TopArtist(ctx, token),
	)
	fmt.Println("1")

	// Retrieve the user's profile data
	// Note: GetUserData does not return an error.
	userProfileInfo := GetUserData(ctx, token)
	fmt.Println("2")

	// Create a new playlist for the user. This function returns an error, so we handle it.
	playlistStatus, err := CreatePlaylist(ctx, token, userProfileInfo.SpotifyID, PlayListName, PlaylistDescription)
	handleError(err, "Failed to Generate new Playlist on user's Profile")
	fmt.Println("3")

	// Build a new music playlist from the returned status
	newPlaylist := NewMusicPlaylist(playlistStatus.URI)
	fmt.Println("4")

	// Construct the final DB document (UserMongoDocument)
	// combining user profile data, music information, and the playlist.
	// NOTE: You may also want to store the playlist ID (playlistStatus.ID) if needed.
	return NewDBDocument(*userProfileInfo, *userMusicInfo, *newPlaylist), nil
}
