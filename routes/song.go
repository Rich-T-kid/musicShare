package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Rich-T-kid/musicShare/pkg/models"
	client "github.com/Rich-T-kid/musicShare/reccommendations/grpc"
	sw "github.com/Rich-T-kid/musicShare/spotwrapper"
)

func GetSongRecommendation(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the path parameters
	vars := mux.Vars(r)
	userID, exists := vars["userID"]
	if !exists || userID == "" {
		logger.Info("Songs endpoint: missing userID in path parameters")
		http.Error(w, "UserID parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch recommendations
	ctx := r.Context()
	songs, err := client.GetReccomendations(ctx, userID)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error generating 'New Song of the day' for user %s: %v", userID, err))
		http.Error(w, "Internal server error while generating a new song", http.StatusInternalServerError)
		return
	}

	// Ensure songs exist in the database
	for _, song := range songs {
		if err := sw.AddSongtoDB(song.SongUri); err != nil {
			logger.Warning(fmt.Sprintf("Failed to ensure song exists: %s - %v", song.SongUri, err))
		}
	}

	// Return songs
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}
func AddSongToDatabase(w http.ResponseWriter, r *http.Request) {
	// Read request body
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("AddSongToDatabase: error reading request body: %v", err))
		http.Error(w, "Malformed JSON body", http.StatusBadRequest)
		return
	}

	// Parse JSON into SongTypes struct
	var song models.SongTypes
	err = json.Unmarshal(bodyByte, &song)
	if err != nil {
		logger.Info(fmt.Sprintf("AddSongToDatabase: error parsing JSON body: %v", err))
		http.Error(w, "Malformed JSON body", http.StatusBadRequest)
		return
	}

	if song.SongURI == "" {
		logger.Info("AddSongToDatabase: missing SongURI field")
		http.Error(w, "SongURI is required", http.StatusBadRequest)
		return
	}

	// Add song to the database
	if err := sw.AddSongtoDB(song.SongURI); err != nil {
		logger.Warning(fmt.Sprintf("Failed to insert song: %s - %v", song.SongURI, err))
		http.Error(w, "Failed to insert song into database", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Successfully added song: %s", song.SongURI)))
}
func GetSongByID(w http.ResponseWriter, r *http.Request) {
	// Extract songID from path parameters
	vars := mux.Vars(r)
	SongUri, exists := vars["songID"]
	if !exists || SongUri == "" {
		logger.Info("GetSongByID: missing songID in path parameters")
		http.Error(w, "SongID parameter is required", http.StatusBadRequest)
		return
	}

	// Query the database for the song
	song, err := sw.ReturnSongbyID(SongUri)
	if err != nil {
		logger.Warning(fmt.Sprintf("GetSongByID: Failed to find song %s - %v", SongUri, err))
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Return the song data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

// tested and works
func Comments(w http.ResponseWriter, r *http.Request) {
	// request Body contains the song id
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Could not read request body. Ensure input body matches api spec")))
			logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
			return
		}
		var request models.CommentsRequest
		err = json.Unmarshal(bodyBytes, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Could not read request body error : %e. Ensure input body matches api spec", err)))
			logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
			return
		}
		// submit Comment under a song
		// TODO: This needs to return the comment UUID to the client
		newUUID, err := sw.SubmitComment(request.SongURI, request.UserResp)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error has occured submiting comment")))
			return
		}
		w.WriteHeader(http.StatusOK)
		response := struct {
			NEW_UUID string `json:"comment_uuid"`
		}{
			NEW_UUID: newUUID,
		}
		json.NewEncoder(w).Encode(response)
	case "GET":
		songURI := r.URL.Query().Get("songURI")
		if songURI == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required query parameter: songURI"))
			return
		}

		comments, err := sw.GetComments(songURI, 0, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error fetching comments"))
			return
		}

		w.WriteHeader(http.StatusOK)
		if len(comments) == 0 {
			json.NewEncoder(w).Encode([]string{})
		}
		json.NewEncoder(w).Encode(comments)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
		return
	}

}

// Tested and works
func CommentsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		comment, err := sw.GetComment(commentID)
		if err != nil { // doesnt exist
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("no comment exist with the commentID passed in %s", commentID)))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
	case "PUT":
		var newComment models.UserComments
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not read request body. Ensure input body matches api spec"))
			logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
			return

		}
		err = json.Unmarshal(bodyBytes, &newComment)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Could not read request body error : %e. Ensure input body matches api spec", err)))
			logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
			return
		}

		found, err := sw.UpdateComment(commentID, newComment)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error occured while trying to update comment %e", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		if found {
			w.Write([]byte(fmt.Sprintf("Found comment with id %s and updated it to %v", commentID, newComment)))
			return
		}
		w.Write([]byte(fmt.Sprintf("Could not find comment with id %s", commentID)))
		return
	case "DELETE":
		err := sw.DeleteComment(commentID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info(fmt.Sprintf("Error occured attempting to delete comment with id %s error: %e", commentID, err))
			return
		}
		w.WriteHeader(200)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
		return
	}

}

// works
func UserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		UserDoc, err := sw.GetUserDocument(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Ensure that a valid userID is pass into the url. %s resulted in this error: %e", userID, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserDoc)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
		return
	}

}

// doesnt work
func UserSongs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		SongTypes, err := sw.GetUserSongs(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Ensure that a valid userID is pass into the url. %s resulted in this error: %e", userID, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SongTypes)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
		return

	}
}

// works
func UserComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	switch r.Method {
	case "GET":
		Comments, err := sw.GetUserComments(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Ensure that a valid userID is pass into the url. %s resulted in this error: %e", userID, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Comments)
		fmt.Println("")
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
		return
	}
}
