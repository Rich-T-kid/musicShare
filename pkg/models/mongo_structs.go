package models

import "time"

type UserMongoDocument struct {
	UUID                string              `json:"ID" bson:"uuid"`
	UserProfileResponse UserProfileResponse `json:"user_profile_response" bson:"user_profile_response"`
	UserMusicInfo       UserMusicInfo       `json:"user_music_info" bson:"user_music_info"`
	MusicSharePlaylist  MusicSharePlaylist  `json:"music_share_playlist" bson:"music_share_playlist"`
	Comments            []UserComments      `json:"use_comments" bson:"use_comments"`
	LikedSongs          []SpotifyURI        `json:"liked_songs" bson:"liked_songs"`
	DislikedSongs       []SpotifyURI        `json:"disliked_songs" bson:"disliked_songs"`
	Listened            []SpotifyURI        `json:"listened" bson:"listened"`
	CreatedAt           time.Time           `json:"created_at" bson:"created_at"`
	Updated             time.Time           `json:"updated" bson:"updated"`
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
type UserTopArtist struct {
	Name        string   `json:"name"`
	URI         string   `json:"uri"`
	Genres      []string `json:"genres"`
	ArtistPhoto string   `json:"ArtistPhoto"`
}
