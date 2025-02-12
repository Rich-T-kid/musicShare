package models

type SignIn struct {
	Username string `json:"username"`
}

// Only being used to pass username through context.context throughout the life cycle of the request
type UsernameKey struct{}

/*
Whenever the application makes requests related to authentication or authorization to Web API,
such as retrieving an access token or refreshing an access token, the error response follows RFC 6749 on the OAuth 2.0 Authorization Framework.
*/
type SpotifyAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/*
Regular Error Object
Apart from the response code, unsuccessful responses return a JSON object containing the following information
*/
type SpotifyError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Exp         int    `json:"exp"`
	Refresh     string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	Expir       int    `json:"exp"`
}

type UserResponse struct {
	Items []struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsPlayable           bool   `json:"is_playable"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		DiscNumber       int      `json:"disc_number"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
		ExternalIds      struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string      `json:"href"`
		ID          string      `json:"id"`
		IsLocal     bool        `json:"is_local"`
		IsPlayable  bool        `json:"is_playable"`
		Name        string      `json:"name"`
		Popularity  int         `json:"popularity"`
		PreviewURL  interface{} `json:"preview_url"`
		TrackNumber int         `json:"track_number"`
		Type        string      `json:"type"`
		URI         string      `json:"uri"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
}

// TODO:Delete
func (u UserTopArtist) top() {}
func (u UserResponse) top()  {}

type SpotifyTopResponse interface {
	top()
}
