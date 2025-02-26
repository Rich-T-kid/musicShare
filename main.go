package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Rich-T-kid/musicShare/routes"
	sw "github.com/Rich-T-kid/musicShare/spotwrapper"
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
	db := sw.CreateNewMongoInstance()
	fmt.Println("Response of mongoDB connection function -> ", db.Connected(context.TODO()))
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}
