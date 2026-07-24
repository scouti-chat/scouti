// Package api makes authenticated calls to the FounderPing /api/v1 gate. The gate
// resolves the access key to a user and enforces row-level security, so the CLI
// only needs to attach the bearer token and relay method/path/body.
package api

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/founderping/founderping/cli/internal/config"
)

// Response is a decoded HTTP response: raw body plus the status code.
type Response struct {
	Status int
	Body   []byte
}

// OK reports whether the status is 2xx.
func (r *Response) OK() bool { return r.Status >= 200 && r.Status < 300 }

// Do performs one authenticated request against /api/v1<path>. body may be nil.
func Do(method, path, accessKey string, body []byte) (*Response, error) {
	url := config.APIBaseURL() + "/api/v1" + normalizePath(path)

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(strings.ToUpper(method), url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Status: res.StatusCode, Body: data}, nil
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}
