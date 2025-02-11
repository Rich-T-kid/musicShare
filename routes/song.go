package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// placing the Crud of route request Now
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

// tested and works
func Comments(w http.ResponseWriter, r *http.Request) {
	// request Body contains the song id
	w.Header().Set("Content-Type", "application/json")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Could not read request body. Ensure input body matches api spec")))
		logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
		return
	}
	var request sw.CommentsRequest
	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Could not read request body error : %e. Ensure input body matches api spec", err)))
		logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
		return
	}
	switch r.Method {
	case "POST":
		// submit Comment under a song
		err = sw.SubmitComment(request.SongURI, request.UserResp)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("")))
			return
		}
		w.WriteHeader(http.StatusOK)

		fmt.Println("")
	case "GET":
		// returns all comments associated with a songURI that we have
		// for now just return all but later add standard api practices like limiting and offsets, ect
		comments, err := sw.GetComments(request.SongURI, 0, 0)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("")))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comments)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("")))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
	case "PUT":
		var newComment sw.UserComments
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
		}
		w.WriteHeader(http.StatusOK)
		if found {
			w.Write([]byte(fmt.Sprintf("Found comment with id %s and updated it to %v", commentID, newComment)))
			return
		}
		w.Write([]byte(fmt.Sprintf("Could not find comment with id %s", commentID)))
	case "DELETE":
		err := sw.DeleteComment(commentID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info(fmt.Sprintf("Error occured attempting to delete comment with id %s error: %e", commentID, err))
		}
		w.WriteHeader(200)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
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
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserDoc)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
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
			fmt.Println("Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Ensure that a valid userID is pass into the url. %s resulted in this error: %e", userID, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SongTypes)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))

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
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Comments)
		fmt.Println("")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
	}
}
