package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/Rich-T-kid/musicShare/routes"
	sw "github.com/Rich-T-kid/musicShare/spotwrapper"
)

var (
	port = "8080"
)

// TODO: make is so that songs are added to the users JSON blob as well when they submit a comment or get a song Review.

func startGRPCServer() {
	cmd := exec.Command("bash", "-c", "source reccommendations/grpc/venv/bin/activate && python3 reccommendations/grpc/server.py") // Example: list files in long format
	res, err := cmd.Output()
	if err != nil {
		log.Fatal("Python Sever Failed to start with an error of: ", err)
		return
	}
	fmt.Printf("Python GRPC server Response: %s\n", string(res))
}
func init() {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	db := sw.CreateNewMongoInstance()
	if db.Connected(context.TODO()) != nil {
		log.Fatal("MongoDB is not connected")
	}
	_ = sw.NewCache[string, string]()
	go startGRPCServer()
	fmt.Print("All external Infra is good to Go!! \n \n")
}

func main() {

	r := routes.InitRoutes() // /exist/r/ == /exist/r
	addr := fmt.Sprintf(":%s", port)
	fullSource := fmt.Sprintf("http://localhost:%s/test", port)
	fmt.Println("Server is running on port", addr, " \n Full url: ", fullSource)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Server has been Ended by error :-> ", err)
	}

}
