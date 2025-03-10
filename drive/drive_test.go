package drive_test

import (
	"testing"

	"go-gdrive-sync/drive"

	"github.com/stretchr/testify/assert"
)

const (
	testTokenFile       = "/home/arunkhattri/github/akNG71/go-forge/go-gdrive-sync/token.json"
	testCredentialsFile = "/home/arunkhattri/github/akNG71/go-forge/go-gdrive-sync/credential.json"
)

func TestListFiles(t *testing.T) {
	// Initialize Drive Service
	srv, err := drive.InitializeDriveService(testTokenFile, testCredentialsFile)
	assert.NoError(t, err, "Drive service should initialize without error")

	// Get list of files
	files, err := drive.ListFiles(srv, 5)
	assert.NoError(t, err, "ListFiles should execute without error")

	// Check if files are retrieved
	assert.GreaterOrEqual(t, len(files), 0, "File list should be non-empty or empty without error")
}
