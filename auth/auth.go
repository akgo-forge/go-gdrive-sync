package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	tokenFile  = "token.json"
	driveScope = "https://www.googleapis.com/auth/drive"
)

// load credentials and return an Oauth2 config with google drive scope
func GetOAuthConfig(credentialFilePath string) (*oauth2.Config, error) {
	data, err := os.ReadFile(credentialFilePath)
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
	token, err := GetTokenFromFile(tokenFile)
	if err != nil {
		// no valid token, request a new one
		token, err = AuthenticateUser(config)
		if err != nil {
			return nil, err
		}
		SaveToken(tokenFile, token)
	}
	return config.Client(context.Background(), token), nil
}
func AuthenticateUser(config *oauth2.Config) (*oauth2.Token, error) {
	// Generate an auth URL
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Open this URL in your browser and authenticate the app:")
	fmt.Println(url)

	// Attempt to open the browser automatically
	if err := openBrowser(url); err != nil {
		fmt.Println("Please open the link manually.")
	}

	// Read the auth code from user input
	fmt.Print("Enter the authentication code: ")
	reader := bufio.NewReader(os.Stdin)
	authCode, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %v", err)
	}
	authCode = strings.TrimSpace(authCode)

	// Exchange the auth code for an access token
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code for token: %v", err)
	}

	// Save the token for future use
	err = SaveToken("token.json", token)
	if err != nil {
		return nil, fmt.Errorf("failed to save token: %v", err)
	}

	return token, nil
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

func GetTokenFromFile(filepath string) (*oauth2.Token, error) {
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

func SaveToken(filepath string, token *oauth2.Token) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(token)
}
