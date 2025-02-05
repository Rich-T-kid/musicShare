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

/*
TODO: move these routes to the / routes folder and keep the main file minimal.
	  Define what information we want user documents to contain for mongoDB (username,email,top songs, top albums, ect)
	  Set up functions and endpoints to handle this
	  cache all of the data the spotifys endpoints return so we dont need to keep calling spotify
	  @
	  Set up CDN  (AWS -> tyler .B and tyler .S)
	  upload to S3 instances and test all endpoints and routes using curl locally as well as postman over the network
	  intergrate with front end
	  Release
*/

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
