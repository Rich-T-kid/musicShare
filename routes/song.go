package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	sw "loveShare/spotWrapper"
)

// placing the Crud of route request Now
func Song(w http.ResponseWriter, r *http.Request) {
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

// tested and works
func Comments(w http.ResponseWriter, r *http.Request) {
	// request Body contains the song id
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
	switch r.Method {
	case "GET":
		SongTypes, err := sw.GetUserSongs(userID)
		if err != nil {
			fmt.Println("Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Ensure that a valid userID is pass into the url. %s resulted in this error: %e", userID, err)))
			return
		}
		fmt.Println("got to here songTypes,", SongTypes)
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
