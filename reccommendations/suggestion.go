package reccommendations

import "fmt"

// hard code for now just for testing purposes
func NewSong(username string, exlude []string) (string, error) {
	fmt.Printf("UserName passed in %s , excluded songs %v", username, exlude)
	return "spotify:track:79rneIpoGKwhqE1MSaw4Ls", nil
}
