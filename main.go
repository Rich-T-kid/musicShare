package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Rich-T-kid/musicShare/routes"
)

var (
	port = "80"
)

/*
API Walkthrough & Validation Testing
(1)
(2)
Define how the front-end authenticates users with the back-end.
Decide on authentication method (tokens, username/password, etc.).
Handle duplicate usernames and define security best practices.
*/

func main() {
	fmt.Printf("mongoDB connection uri %s\n redis connection string %s\n ", os.Getenv("MONGO_URI"), os.Getenv("REDIS_ADDR"))
	r := routes.InitRoutes() // /exist/r/ == /exist/r

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}
