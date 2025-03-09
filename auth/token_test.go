package auth_test

import (
	"os"
	"testing"

	"go-gdrive-sync/auth"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

const testTokenFile = "test_token.json"

func TestSaveAndGetToken(t *testing.T) {
	// Create a sample token
	expectedToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
	}

	// Save token to file
	err := auth.SaveToken(testTokenFile, expectedToken)
	assert.NoError(t, err, "SaveToken should not return an error")

	// Retrieve token from file
	retrievedToken, err := auth.GetTokenFromFile(testTokenFile)
	assert.NoError(t, err, "GetTokenFromFile should not return an error")
	assert.Equal(t, expectedToken.AccessToken, retrievedToken.AccessToken, "AccessToken should match")
	assert.Equal(t, expectedToken.RefreshToken, retrievedToken.RefreshToken, "RefreshToken should match")

	// Cleanup test file
	os.Remove(testTokenFile)
}
