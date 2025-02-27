package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Printf("mongoDB connection uri %s\n redis connection string %s\n ", os.Getenv("MONGO_URI"), os.Getenv("REDIS_ADDR"))
	r := routes.InitRoutes() // /exist/r/ == /exist/r
	db := sw.CreateNewMongoInstance()
	fmt.Println("Response of mongoDB connection function -> ", db.Connected(context.TODO()))
	if db.Connected(context.TODO()) != nil {
		log.Fatal("MongoDB is not connected")
	}
	cache := sw.NewCache[string, string]()
	cache.Set(context.TODO(), "Richard", "king", 10)
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port", addr)
	http.ListenAndServe(addr, r)

}
