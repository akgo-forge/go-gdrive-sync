package gdrive_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"

	"go-gdrive-sync/gdrive"
)

const (
	testTokenFile       = "/home/arunkhattri/github/akNG71/go-forge/go-gdrive-sync/token.json"
	testCredentialsFile = "/home/arunkhattri/github/akNG71/go-forge/go-gdrive-sync/credential.json"
)

func TestListFiles(t *testing.T) {
	// Initialize Drive Service
	srv, err := gdrive.InitializeDriveService(testTokenFile, testCredentialsFile)
	assert.NoError(t, err, "Drive service should initialize without error")

	// Get list of files
	files, err := gdrive.ListFiles(srv, 5)
	assert.NoError(t, err, "ListFiles should execute without error")

	// Check if files are retrieved
	assert.GreaterOrEqual(t, len(files), 0, "File list should be non-empty or empty without error")
}

// MockDriveService mocks the DriveService interface
type MockDriveService struct {
	UploadError bool
}

var _ gdrive.DriveService = (*MockDriveService)(nil)

// Mock the FilesCreate method
func (m *MockDriveService) FilesCreate(file *drive.File, media io.Reader) (*drive.File, error) {
	if m.UploadError {
		return nil, errors.New("mock upload failure")
	}
	return &drive.File{Name: file.Name, Id: "mock-file-id"}, nil
}

func TestUploadFile(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "testfile.txt")
	assert.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tempFile.Name())

	// Write some data to the temp file
	_, err = tempFile.Write([]byte("This is a test file"))
	assert.NoError(t, err, "Failed to write to temporary file")
	tempFile.Close()

	// Extract the base filename
	expectedFileName := filepath.Base(tempFile.Name())

	// Mock Google Drive service
	mockService := &MockDriveService{UploadError: false}

	// Call UploadFile
	uploadedFile, err := gdrive.UploadFile(mockService, tempFile.Name())

	// Verify expectations
	assert.NoError(t, err, "UploadFile should not return an error")
	assert.NotNil(t, uploadedFile, "Uploaded file should not be nil")
	assert.Equal(t, expectedFileName, filepath.Base(uploadedFile.Name), "Uploaded file name should match the original")

	// Test failure scenario
	mockFailService := &MockDriveService{UploadError: true}

	_, err = gdrive.UploadFile(mockFailService, tempFile.Name())
	assert.Error(t, err, "Expected an error when file upload fails")
}
