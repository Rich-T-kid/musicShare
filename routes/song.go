package routes

import (
	"fmt"
	"net/http"
)

// placing the Crud of route request Now
func Song(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Println("")
	case "GET":
		fmt.Println("")
	case "PUT":
		fmt.Println("")
	case "DELETE":
		fmt.Println("")
	}
}

func Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Println("")
	case "GET":
		fmt.Println("")
	case "PUT":
		fmt.Println("")
	case "DELETE":
		fmt.Println("")
	}

}
