package main

import (
	"fmt"
	"net/http"

	"loveShare/routes"
)

var (
	port = "8080"
)

/*
TODO: move these routes to the / routes folder and keep the main file minimal.
	  cache all of the data the spotifys endpoints return so we dont need to keep calling spotify -> half done, Not 100 % sure what endpoints are super critcal right now
	  @
	  Set up CDN  (AWS -> tyler .B and tyler .S)
	  upload to S3 instances and test all endpoints and routes using curl locally as well as postman over the network
	  intergrate with front end
	  Release
*/

func main() {
	r := routes.InitRoutes() // /exist/r/ == /exist/r

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}
