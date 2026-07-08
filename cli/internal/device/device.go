// Package device implements the OAuth-style device authorization flow used by
// `scouti login`. The CLI starts a request, the developer approves it in a
// browser, and the CLI polls with a private device_code to receive the key.
// The high-entropy device_code never leaves the machine; the short user_code is
// only shown so the developer can eyeball-match it on the approval page.
package device

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/scouti-chat/scouti/cli/internal/config"
)

// Start describes an initiated authorization (from POST /auth/device/start).
type Start struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	Interval                int    `json:"interval"`
	ExpiresIn               int    `json:"expires_in"`
}

// PollResult is the outcome of a single poll. When Issued is true the key is
// ready; otherwise wait Interval seconds and poll again.
type PollResult struct {
	Issued    bool
	AccessKey string
	UserID    string
	Interval  int
}

// Begin initiates a device authorization.
func Begin() (*Start, error) {
	res, err := http.Post(endpoint("start"), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("could not start login (%d)", res.StatusCode)
	}
	var s Start
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}
	if s.Interval <= 0 {
		s.Interval = 5
	}
	return &s, nil
}

// Poll checks the authorization once. A terminal condition (expired / denied)
// returns an error; a still-pending state returns Issued=false with no error.
func Poll(deviceCode string) (PollResult, error) {
	payload, _ := json.Marshal(map[string]string{"device_code": deviceCode})
	res, err := http.Post(endpoint("token"), "application/json", bytes.NewReader(payload))
	if err != nil {
		return PollResult{}, err
	}
	defer res.Body.Close()

	var body struct {
		AccessKey string `json:"access_key"`
		UserID    string `json:"user_id"`
		Interval  int    `json:"interval"`
		Error     struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)

	if res.StatusCode/100 == 2 {
		if body.AccessKey != "" {
			return PollResult{Issued: true, AccessKey: body.AccessKey, UserID: body.UserID}, nil
		}
		interval := body.Interval
		if interval <= 0 {
			interval = 5
		}
		return PollResult{Interval: interval}, nil
	}

	msg := body.Error.Message
	if msg == "" {
		msg = fmt.Sprintf("login failed (%d)", res.StatusCode)
	}
	return PollResult{}, errors.New(msg)
}

// OpenBrowser makes a best-effort attempt to open target in the default browser
// and reports whether the launcher started. On headless machines it returns
// false, and the caller falls back to printing the URL.
func OpenBrowser(target string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	if err := cmd.Start(); err != nil {
		return false
	}
	go func() { _ = cmd.Wait() }() // reap without blocking
	return true
}

func endpoint(leaf string) string {
	return config.APIBaseURL() + "/api/v1/auth/device/" + leaf
}
