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
