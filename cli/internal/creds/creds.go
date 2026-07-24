// Package creds persists and resolves the developer's FounderPing access key (UAK).
// The key is secret: it is written owner-only and never printed by the CLI.
package creds

import (
	"encoding/json"
	"os"

	"github.com/founderping/founderping/cli/internal/config"
)

// Credentials is the on-disk shape at ~/.founderping/credentials.json.
type Credentials struct {
	// AccessKey is the user access key (uak_...). It acts as a long-lived
	// refresh credential; the gate exchanges it for short-lived tokens.
	AccessKey string `json:"access_key"`
	UserID    string `json:"user_id,omitempty"`
}

// Read returns the stored credentials, or (nil, nil) if none exist yet.
func Read() (*Credentials, error) {
	path, err := config.CredentialsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var c Credentials
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Write persists credentials, creating ~/.founderping (0700) and the file (0600) so
// the key stays readable only by its owner.
func Write(c Credentials) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	path, err := config.CredentialsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return err
	}
	// WriteFile won't tighten perms on an existing file; enforce 0600.
	return os.Chmod(path, 0o600)
}

// ResolveAccessKey returns the key to authenticate with. FOUNDERPING_ACCESS_KEY (for
// CI) takes precedence over the stored file. Empty string means "not logged in".
func ResolveAccessKey() (string, error) {
	if k := os.Getenv("FOUNDERPING_ACCESS_KEY"); k != "" {
		return k, nil
	}
	c, err := Read()
	if err != nil || c == nil {
		return "", err
	}
	return c.AccessKey, nil
}
