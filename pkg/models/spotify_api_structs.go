package models

import "time"

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
type SpotifyURI struct {
	Song string    `json:"songID`
	Date time.Time `json:"date"`
}

type MusicSharePlaylist struct {
	Name        string       `json:"name"`
	PlaylistURI string       `json:"playlist_uri"`
	Songs       []SpotifyURI `json:"songs"` // filled with
}
