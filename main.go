package main

import (
	"fmt"
	"net/http"

	"github.com/Rich-T-kid/musicShare/routes"
)

var (
	port = "8080"
)

/*
API Walkthrough & Validation Testing
(1)
Test API endpoints using Swagger (front-end perspective).
Manually test API requests with valid & malformed input (user perspective).
Ensure validation works correctly and tokens are stored in MongoDB & cache.
Frontend Authentication & Security
(2)
Define how the front-end authenticates users with the back-end.
Decide on authentication method (tokens, username/password, etc.).
Handle duplicate usernames and define security best practices.
*/

func main() {
	r := routes.InitRoutes() // /exist/r/ == /exist/r

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}

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
