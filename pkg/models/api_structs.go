package models

type CommentsRequest struct {
	SongURI  string       `json:"songID"`
	UserResp UserComments `json:"userComment"`
}

type SongTypes struct {
	SongURI       string         `json:"songURI" bson:"songURI"`
	Comments      []UserComments `json:"comments" bson:"comments"`
	AlternateName []string       `json:"alternateName" bson:"alternateName"`
	UUID          string         `json:"ID" bson:"uuid"`
}

type UserComments struct {
	Username string `json:"username" bson:"username"`
	Rating   uint8  `json:"rating" bson:"rating"` // out of 5
	Review   string `json:"review" bson:"review"`
	SongID   string `json:"songID" bson:"songID"`
	UUID     string `json:"ID" bson:"uuid"`
}
