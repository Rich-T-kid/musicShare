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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Printf("mongoDB connection uri %s\n redis connection string %s redis Password %s \n ", os.Getenv("MONGO_URI"), os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	db := sw.CreateNewMongoInstance()
	if db.Connected(context.TODO()) != nil {
		log.Fatal("MongoDB is not connected")
	}
	_ = sw.NewCache[string, string]()
	fmt.Println("All external Infra is good to Go -> Starting Main function \n \n")
}

func main() {

	r := routes.InitRoutes() // /exist/r/ == /exist/r
	addr := fmt.Sprintf(":%s", port)
	fullSource := fmt.Sprintf("http://localhost:%s/test", port)
	fmt.Println("Server is running on port", addr, " \n Full url: ", fullSource)
	http.ListenAndServe(addr, r)

}
