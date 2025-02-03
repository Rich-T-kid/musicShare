package main

import (
	"fmt"
	"net/http"

	"loveShare/routes"
)

var (
	port = "8080"
)

func main() {
	http.HandleFunc("/test", routes.Test)          //api
	http.HandleFunc("/login", routes.HomePage)     //html
	http.HandleFunc("/signIn", routes.SignIn)      // api
	http.HandleFunc("/link", routes.RedirectLink)  // api
	http.HandleFunc("/callback", routes.SongofDay) // html
	http.HandleFunc("/auth", routes.RedirectPage)  //html
	http.HandleFunc("/Songs", routes.RedirectPage) //api
	Addr := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port ", Addr)
	http.ListenAndServe(Addr, nil)
}
