package spotwrapper

type SignIn struct {
	Username string `json:"username"`
}

type UsernameKey struct{}

/*
Authentication Error Object
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

func (u UserTopArtist) top() {}
func (u UserResponse) top()  {}

type SpotifyTopResponse interface {
	top()
}

type PlaylistResponse struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Images       []interface{} `json:"images"`
	PrimaryColor interface{}   `json:"primary_color"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	URI          string        `json:"uri"`
	Owner        struct {
		Href         string      `json:"href"`
		ID           string      `json:"id"`
		Type         string      `json:"type"`
		URI          string      `json:"uri"`
		DisplayName  interface{} `json:"display_name"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"owner"`
	Public     bool   `json:"public"`
	SnapshotID string `json:"snapshot_id"`
	Tracks     struct {
		Limit    int           `json:"limit"`
		Next     interface{}   `json:"next"`
		Offset   int           `json:"offset"`
		Previous interface{}   `json:"previous"`
		Href     string        `json:"href"`
		Total    int           `json:"total"`
		Items    []interface{} `json:"items"`
	} `json:"tracks"`
}
type CreatePlaylistRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}
type PlayListaddtion struct {
	SnapShot string `json:"snapshot_id"`
}
type PlaylistAddtionRequest struct {
	URI      []string `json:"uris"`
	Position int      `json:"position"`
}
type UserTopArtist struct {
	Name        string   `json:"name"`
	URI         string   `json:"uri"`
	Genres      []string `json:"genres"`
	ArtistPhoto string   `json:"ArtistPhoto"`
}

/*
implementfunctions to parse each needed response so that we can store in mongoDB
*/
type UserProfileResponse struct {
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Images      []struct {
		//Height int    `json:"height"`
		URL string `json:"url"`
		//Width  int    `json:"width"`
	} `json:"images"`
	//ExplicitContent struct {
	//	FilterEnabled bool `json:"filter_enabled"`
	//	FilterLocked  bool `json:"filter_locked"`
	//} `json:"explicit_content"`
	//ExternalUrls struct {
	//	Spotify string `json:"spotify"`
	//} `json:"external_urls"`
	//Followers struct {
	//	Href  interface{} `json:"href"`
	//	Total int         `json:"total"`
	//} `json:"followers"`
	//Href    string        `json:"href"`
	SpotifyID string `json:"id"`
	///Images  []interface{} `json:"images"`
	//Product string        `json:"product"`
	//Type    string        `json:"type"`
	//URI     string        `json:"uri"`
}
type FollowedArtist struct {
	Name    string   `json:"name"`
	Spotify string   `json:"Spotify"`
	Genres  []string `json:"genres"`
	URI     string   `json:"uri"`
}
type Album struct {
	Artist      string `json:"Artist"`
	Name        string `json:"Name"`
	AlbumLink   string `json:"AlbumLink"`
	AlbumURI    string `json:"AlbumURI"`
	AlbumID     string `json:"AlbumID"`
	AlbumImage  Image  `json:"AlbumImage"`
	AlbumName   string `json:"AlbumName"`
	TotalTracks int    `json:"totalTracks"`
	ReleaseDate string `json:"release_date"`
}

type Image struct {
	URL string `json:"url"`
}
type UserMusicInfo struct {
	FollowedArtist []FollowedArtist `json:"FollowedArtist"` // finsihed
	TopTracks      UserTopTrack     `json:"TopTracks"`      // Finished
	TopsArtist     []UserTopArtist  `json:"TopsArtist"`     // Finished
}
type UserTopTrack struct {
	TopAlbums  []Album       `json:"TopAlbums"`
	TopSingles []SingleTrack `json:"TopSinglesTracks"`
}
type SingleTrack struct {
	Artist      string `json:"Artist"`
	Name        string `json:"Name"`
	TrackLink   string `json:"trackLink"`
	TrackName   string `json:"TrackName"`
	ReleaseDate string `json:"release_date"`
}
type TopsArtist struct {
	Name        string   `json:"name"`
	URI         string   `json:"uri"`
	Genres      []string `json:"genres"`
	ArtistPhoto string   `json:"ArtistPhoto"`
}

type SpotArtist struct {
	Artists struct {
		Href    string `json:"href"`
		Limit   int    `json:"limit"`
		Next    string `json:"next"`
		Cursors struct {
			After string `json:"after"`
		} `json:"cursors"`
		Total int `json:"total"`
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  interface{} `json:"href"`
				Total int         `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
			URI        string `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
}
type SpotTracks struct {
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
type SpotTopArtist struct {
	Items []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  interface{} `json:"href"`
			Total int         `json:"total"`
		} `json:"followers"`
		Genres []string `json:"genres"`
		Href   string   `json:"href"`
		ID     string   `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		Type       string `json:"type"`
		URI        string `json:"uri"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
}

type SpotifyTopArtistResponse struct {
	Items []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  interface{} `json:"href"`
			Total int         `json:"total"`
		} `json:"followers"`
		Genres []string `json:"genres"`
		Href   string   `json:"href"`
		ID     string   `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		Type       string `json:"type"`
		URI        string `json:"uri"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
}

type SpotifyTrackResponse struct {
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
			//ExternalUrls     struct {
			//	Spotify string `json:"spotify"`
			//} `json:"external_urls"`
			//		Href   string `json:"href"`
			//		ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			//		IsPlayable           bool   `json:"is_playable"`
			Name                 string `json:"name"` // album name
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"` // track or album
			URI                  string `json:"uri"`  // uri of track or album
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"` // artist spotify page
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`   // artist ID
			Name string `json:"name"` // artist name
			Type string `json:"type"` // will always be artist
			URI  string `json:"uri"`  // spotify uri for artist
		} `json:"artists"`
		//AvailableMarkets []string `json:"available_markets"`
		//DiscNumber       int      `json:"disc_number"`
		//DurationMs       int      `json:"duration_ms"`
		//Explicit         bool     `json:"explicit"`
		//ExternalIds      struct {
		//	Isrc string `json:"isrc"`
		//} `json:"external_ids"`
		//ExternalUrls struct {
		//	Spotify string `json:"spotify"`
		//} `json:"external_urls"`
		//Href        string      `json:"href"`
		ID string `json:"id"`
		//IsLocal     bool        `json:"is_local"`
		//IsPlayable  bool        `json:"is_playable"`
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
