package credentials

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "linear-cli"

// Credential keys
const (
	LinearAPIKey  = "LINEAR_API_KEY"
	GoogleAPIKey  = "GOOGLE_API_KEY"
	LinearTeamID  = "LINEAR_TEAM_ID"
)

// Set stores a credential in the system keyring
func Set(key, value string) error {
	if value == "" {
		return fmt.Errorf("credential value cannot be empty")
	}
	return keyring.Set(serviceName, key, value)
}

// Get retrieves a credential from the system keyring
func Get(key string) (string, error) {
	value, err := keyring.Get(serviceName, key)
	if err != nil {
		return "", fmt.Errorf("failed to get credential %s: %w", key, err)
	}
	return value, nil
}

// Delete removes a credential from the system keyring
func Delete(key string) error {
	return keyring.Delete(serviceName, key)
}

// GetLinearAPIKey retrieves the Linear API key
func GetLinearAPIKey() (string, error) {
	return Get(LinearAPIKey)
}

// GetGoogleAPIKey retrieves the Google API key
func GetGoogleAPIKey() (string, error) {
	return Get(GoogleAPIKey)
}

// GetLinearTeamID retrieves the Linear team ID
func GetLinearTeamID() (string, error) {
	return Get(LinearTeamID)
}

// SetLinearAPIKey stores the Linear API key
func SetLinearAPIKey(value string) error {
	return Set(LinearAPIKey, value)
}

// SetGoogleAPIKey stores the Google API key
func SetGoogleAPIKey(value string) error {
	return Set(GoogleAPIKey, value)
}

// SetLinearTeamID stores the Linear team ID
func SetLinearTeamID(value string) error {
	return Set(LinearTeamID, value)
}
