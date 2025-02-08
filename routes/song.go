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
		w.Write([]byte(fmt.Sprintf("Could not read request body. Ensure input body matches api spec")))
		logger.Warning(fmt.Sprintf("Error decoding requst body %e", err))
		return
	}
	switch r.Method {
	case "POST":
		fmt.Println("")
	case "GET": // returns all comments associated with a songURI that we have
		// for now just return all but later add standard api practices like limiting and offsets, ect
		fmt.Println("")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
	}

}

func CommentsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["comment_id"]
	fmt.Println(userID)
	switch r.Method {
	case "GET":
		fmt.Println("")
	case "PUT":
		fmt.Println("")
	case "DELETE":
		fmt.Println("")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
	}

}

func UserID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
	}

}
func UserSongs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))

	}
}
func UserComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed", r.Method)))
	}
}
