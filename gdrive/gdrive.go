package gdrive

import (
	"context"
	"fmt"
	"go-gdrive-sync/auth"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveService interface {
	FilesCreate(file *drive.File, media io.Reader) (*drive.File, error)
}

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

// file mime type
func detectMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the first 512 byte to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// detect MIME type
	mimeType := http.DetectContentType(buffer)
	if mimeType == "application/octet-stream" {
		// generic, get from extension
		extMimeType := mime.TypeByExtension(filepath.Ext(filePath))
		if extMimeType != "" {
			return extMimeType, nil
		}
	}
	return mimeType, nil
}

// Upload files
func UploadFile(service DriveService, filePath string) (*drive.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	driveFile := &drive.File{Name: filePath}
	uploadedFile, err := service.FilesCreate(driveFile, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	return uploadedFile, nil
}

// func UploadFile(service DriveService, filePath string) (*drive.File, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	fileName := filepath.Base(filePath)
// 	// MIME type
// 	mimeType, err := detectMimeType(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to detect MIME type: %v", err)
// 	}

// 	// file metadata
// 	driveFile := &drive.File{
// 		Name: fileName,
// 	}

// 	// upload
// 	fileCall := service.FilesCreate(driveFile).
// 		Media(file, googleapi.ContentType(mimeType))
// 	uploadedFile, err := fileCall.Do()

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to upload file: %v", err)
// 	}
// 	fmt.Printf("Uploaded file '%s' (ID: %s)\n", uploadedFile.Name, uploadedFile.Id)
// 	return uploadedFile, nil
// }
