package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {

	r := mux.NewRouter().StrictSlash(true) // /exist/r/ == /exist/r
	r.Use(enableCors)
	// Fine
	// Define routes
	r.HandleFunc("/test", Test).Methods("GET")         // API
	r.HandleFunc("/link", RedirectLink).Methods("GET") // API
	r.HandleFunc("/callback", Callback).Methods("GET") // HTML

	// not fine
	r.HandleFunc("/login", HomePage).Methods("GET")      // HTML
	r.HandleFunc("/signIn", SignIn).Methods("POST")      // API
	r.HandleFunc("/auth", RedirectPage).Methods("GET")   // HTML
	r.HandleFunc("/loveShare", LoveShare).Methods("GET") // HTML

	// Basic Crud
	r.HandleFunc("/song/{songID}", GetSongByID).Methods("GET")

	r.HandleFunc("/songs/{userID}", GetSongRecommendation).Methods("GET") // Fetch song recommendation
	r.HandleFunc("/songs/add", AddSongToDatabase).Methods("POST")         // Manually add a song

	r.HandleFunc("/comments", Comments).Methods("GET", "POST", "PUT", "DELETE")        // API
	r.HandleFunc("/comments/{comment_id}", CommentsID).Methods("GET", "PUT", "DELETE") // API
	r.HandleFunc("/users/{user_id}", UserID).Methods("GET")                            // return user json document                            // API
	r.HandleFunc("/users/{user_id}/songs", UserSongs).Methods("GET")                   // return all song uri's a user has listneded to as well as their liked and disliked songs                   // API
	r.HandleFunc("/users/{user_id}/comments", UserComments).Methods("GET")             // return all coments made by user by their user id            // API
	// Start server
	return r
}

// TODO: this isnt done finish
// middleware Mabey will be used in auth later
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("url: %s Method: %s Origin: %s\n", r.URL.Path, r.Method, r.Header.Get("Origin"))

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or your specific origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
