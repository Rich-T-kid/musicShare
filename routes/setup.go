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
	r.HandleFunc("/test", Test).Methods("GET")           // API
	r.HandleFunc("/login", HomePage).Methods("GET")      // HTML
	r.HandleFunc("/signIn", SignIn).Methods("POST")      // API
	r.HandleFunc("/link", RedirectLink).Methods("GET")   // API
	r.HandleFunc("/callback", SongofDay).Methods("GET")  // HTML
	r.HandleFunc("/loveShare", LoveShare).Methods("GET") // HTML
	r.HandleFunc("/auth", RedirectPage).Methods("GET")   // HTML
	r.HandleFunc("/Songs", RedirectPage).Methods("GET")  // API
	r.HandleFunc("/api/Songs", Song).Methods("POST")
	r.HandleFunc("/exist/{name}", UniqueUsername).Methods("GET") // API
	// Start server
	return r
}

// Middleware
func temp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Middleware Here just testing some things")
		next.ServeHTTP(w, r)
	})
}
