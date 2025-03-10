package drive

import (
	"context"
	"fmt"
	"go-gdrive-sync/auth"
	"log"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func InitializeDriveService(tokenFile string, credentialFile string) (*drive.Service, error) {
	// Get OAuth2 config
	config, err := auth.GetOAuthConfig(credentialFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth config: %v", err)
	}

	// Get token from file
	token, err := auth.GetTokenFromFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Create OAuth2 client
	client := config.Client(context.Background(), token)

	// initialize Drive Service
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Drive Service: %v", err)
	}

	log.Println("Google Drive client initialized successfully")
	return srv, nil
}

// retrieve list of files from Google Drive
func ListFiles(srv *drive.Service, maxResults int) ([]*drive.File, error) {
	filesList := []*drive.File{}

	// Call the Drive API
	fileList, err := srv.Files.List().
		PageSize(int64(maxResults)).
		Fields("files(id, name, mimeType, size)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	filesList = fileList.Files
	return filesList, nil

}
