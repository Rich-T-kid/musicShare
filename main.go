package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"loveShare/routes"
)

var (
	port = "8080"
)

func main() {
	r := mux.NewRouter().StrictSlash(true) // /exist/r/ == /exist/r

	// Define routes
	r.HandleFunc("/test", routes.Test).Methods("GET")                   // API
	r.HandleFunc("/login", routes.HomePage).Methods("GET")              // HTML
	r.HandleFunc("/signIn", routes.SignIn).Methods("POST")              // API
	r.HandleFunc("/link", routes.RedirectLink).Methods("GET")           // API
	r.HandleFunc("/callback", routes.SongofDay).Methods("GET")          // HTML
	r.HandleFunc("/loveShare", routes.LoveShare).Methods("GET")         // HTML
	r.HandleFunc("/auth", routes.RedirectPage).Methods("GET")           // HTML
	r.HandleFunc("/Songs", routes.RedirectPage).Methods("GET")          // API
	r.HandleFunc("/exist/{name}", routes.UniqueUsername).Methods("GET") // API
	// Start server
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}
