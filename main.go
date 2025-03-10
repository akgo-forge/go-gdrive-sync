package main

import (
	"fmt"
	"go-gdrive-sync/auth"
	"log"
)

func main() {
	config, err := auth.GetOAuthConfig("credential.json")
	if err != nil {
		log.Fatalf("Error loading OAuth config: %v", err)
	}

	token, err := auth.AuthenticateUser(config)
	if err != nil {
		log.Fatalf("Error authenticating user: %v", err)
	}

	fmt.Println("Authentication successful! Token saved.")
	fmt.Printf("Access Token: %s\n", token.AccessToken)
}
