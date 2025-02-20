package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {

	r := mux.NewRouter().StrictSlash(true) // /exist/r/ == /exist/r
	r.Use(temp)
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
	r.HandleFunc("/songs", SongOfTheDay).Methods("POST")                               // API
	r.HandleFunc("/comments", Comments).Methods("GET", "POST", "PUT", "DELETE")        // API
	r.HandleFunc("/comments/{comment_id}", CommentsID).Methods("GET", "PUT", "DELETE") // API
	r.HandleFunc("/users/{user_id}", UserID).Methods("GET")                            // return user json document                            // API
	r.HandleFunc("/users/{user_id}/songs", UserSongs).Methods("GET")                   // return all song uri's a user has listneded to as well as their liked and disliked songs                   // API
	r.HandleFunc("/users/{user_id}/comments", UserComments).Methods("GET")             // return all coments made by user by their user id            // API
	// Start server
	return r
}

// middleware Mabey will be used in auth later
func temp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("url: %s Method: %s ", r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}
