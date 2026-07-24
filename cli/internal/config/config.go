// Package config resolves the API origin and on-disk paths the CLI depends on.
package config

import (
	"os"
	"path/filepath"
	"strings"
)

// defaultAPIURL is the production origin that serves the /api/v1 gate.
const defaultAPIURL = "https://founderping.app"

// dirName is the per-user directory under $HOME that holds credentials.
const dirName = ".founderping"

// APIBaseURL returns the API origin without a trailing slash. FOUNDERPING_API_URL
// overrides it for staging / local development.
func APIBaseURL() string {
	v := os.Getenv("FOUNDERPING_API_URL")
	if v == "" {
		v = defaultAPIURL
	}
	return strings.TrimRight(v, "/")
}

// Dir is the config directory (~/.founderping).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dirName), nil
}

// CredentialsPath is the file that stores the access key (~/.founderping/credentials.json).
func CredentialsPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}
