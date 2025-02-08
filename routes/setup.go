package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {

	r := mux.NewRouter().StrictSlash(true) // /exist/r/ == /exist/r
	r.Use(temp)
	// Define routes
	r.HandleFunc("/test", Test).Methods("GET")                                         // API
	r.HandleFunc("/login", HomePage).Methods("GET")                                    // HTML
	r.HandleFunc("/signIn", SignIn).Methods("POST")                                    // API
	r.HandleFunc("/link", RedirectLink).Methods("GET")                                 // API
	r.HandleFunc("/callback", SongofDay).Methods("GET")                                // HTML
	r.HandleFunc("/loveShare", LoveShare).Methods("GET")                               // HTML
	r.HandleFunc("/auth", RedirectPage).Methods("GET")                                 // HTML
	r.HandleFunc("/songs", RedirectPage).Methods("GET")                                // API
	r.HandleFunc("/comments", Comments).Methods("GET", "POST", "PUT", "DELETE")        // API
	r.HandleFunc("/comments/{comment_id}", CommentsID).Methods("GET", "PUT", "DELETE") // API
	r.HandleFunc("/Users/{user_id}", UserID).Methods("GET")                            // API
	r.HandleFunc("/Users/{user_id}/songs", UserSongs).Methods("GET")                   // API
	r.HandleFunc("/Users/{user_id}/comments", UserComments).Methods("GET")             // API
	r.HandleFunc("/exist/{name}", UniqueUsername).Methods("GET")                       // API
	// Start server
	return r
}

func temp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Middleware Here just testing some things")
		next.ServeHTTP(w, r)
	})
}
