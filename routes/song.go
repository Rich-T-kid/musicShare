package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	rec "loveShare/reccomendations"
)

type SongRequest struct {
	UserName    string   `json:"username"`
	ExcludeList []string `json:"unwanted_tracks"`
}

// placing the Crud of route request Now
// Song of the Day
func Song(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		bodyByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Info(fmt.Sprintf("Songs endpoint: error reading request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed JSON body"))
			return // Early return
		}

		// 2) Parse JSON into SongRequest
		var requestJson SongRequest
		err = json.Unmarshal(bodyByte, &requestJson)
		if err != nil {
			logger.Info(fmt.Sprintf("Songs endpoint: error parsing JSON body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed JSON body"))
			return // Early return
		}
		//cache := userNamecache
		if requestJson.UserName == "" { //|| cache.Exist(r.Context(),fmt.Sprintf("UniqueUserName:%s",requestJson.UserName)){
			logger.Info("Songs endpoint: empty or invalid UserName field")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Username cannot be empty or invalid"))
			return // Early return
		}

		// 4) Business logic: generate a new song
		song, err := rec.NewSong(requestJson.UserName, requestJson.ExcludeList)
		if err != nil {
			logger.Warning(fmt.Sprintf("Error generating 'New Song of the day' for user %s: %v", requestJson.UserName, err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error while generating a new song"))
			return // Early return
		}

		// 5) Return success
		response := map[string]string{
			"SpotifyURI": song,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("%s is not a valid method for this route", r.Method)))
	}
}

func Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Println("")
	case "GET":
		fmt.Println("")
	case "PUT":
		fmt.Println("")
	case "DELETE":
		fmt.Println("")
	}

}
