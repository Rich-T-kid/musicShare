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
