package main

import (
	"fmt"
	"net/http"

	"github.com/Rich-T-kid/musicShare/routes"
)

var (
	port = "80"
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
