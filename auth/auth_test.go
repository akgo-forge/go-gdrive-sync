package auth_test

import (
	"os"
	"testing"

	"go-gdrive-sync/auth"

	"github.com/stretchr/testify/assert"
)

const testCredentialsFile = "test_credentials.json"

func TestGetOAuthConfig(t *testing.T) {
	mockCredentials := `{
		"installed": {
			"client_id": "test-client-id",
			"client_secret": "test-client-secret",
			"redirect_uris": ["http://localhost"]
		}
	}`

	err := os.WriteFile(testCredentialsFile, []byte(mockCredentials), 0644)
	assert.NoError(t, err, "Should create test credentials file successfully")

	// Call function with test credentials file path
	config, err := auth.GetOAuthConfig(testCredentialsFile)
	assert.NoError(t, err, "GetOAuthConfig should not return an error")
	assert.Equal(t, "test-client-id", config.ClientID, "ClientID should match")
	assert.Equal(t, "test-client-secret", config.ClientSecret, "ClientSecret should match")

	// Cleanup
	os.Remove(testCredentialsFile)
}

// NOTE: Testing GetTokenFromWeb is tricky because it requires user interaction.
// You can mock the HTTP response if needed.
