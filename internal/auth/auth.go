package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	credentialFile = "./client_secret_150306401807-tg3007pjnckcdfcr0nr3nvb9qmu1rcar.apps.googleusercontent.com.json"
	tokenFile      = "token.json"
	driveScope     = "https://wwww.googleapis.com/auth/drive"
)

// load credentials and return an Oauth2 config with google drive scope
func GetOAuthConfig() (*oauth2.Config, error) {
	data, err := os.ReadFile(credentialFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(data, driveScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials.json: %v", err)
	}

	return config, nil
}

// retrieves an OAuth2 client using a saved token or authenticate the user
func GetClient(config *oauth2.Config) (*http.Client, error) {
	token, err := getTokenFromFile(tokenFile)
	if err != nil {
		// no valid token, request a new one
		token, err = AuthenticateUser(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokenFile, token)
	}
	return config.Client(context.Background(), token), nil
}

func AuthenticateUser(config *oauth2.Config) (*oauth2.Token, error) {
	// Generate an auth URL and open it in a browser
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Open this URL in your browser and authenticate the app:")
	fmt.Println(url)

	// open browser automatically
	if err := openBrowser(url); err != nil {
		fmt.Println("Please open the link manually.")
	}

	// Read the auth code from usr input
	fmt.Print("Enter the authentication code: ")
	var authCode string
	fmt.Scan(&authCode)

	// Exchange the auth code for an access token
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code for token: %v", err)
	}
	return token, nil
}

func getTokenFromFile(filepath string) (*oauth2.Token, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func saveToken(filepath string, token *oauth2.Token) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(token)
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("Unsupported platform")
	}
}
