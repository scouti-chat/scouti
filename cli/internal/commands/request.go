package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/scouti-chat/scouti/cli/internal/api"
	"github.com/scouti-chat/scouti/cli/internal/creds"
)

// Request forwards one authenticated call to /api/v1 and prints the response
// body. The exit code mirrors the HTTP result (0 for 2xx) so an agent can tell
// success from failure without parsing.
func Request(args []string) int {
	key, err := creds.ResolveAccessKey()
	if err != nil {
		return fail(err)
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "Not logged in. Run `scouti login` first (or set SCOUTI_ACCESS_KEY).")
		return 1
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: scouti request <METHOD> <PATH> [body]")
		fmt.Fprintln(os.Stderr, "  body: path/to/file.json | @file.json | - (stdin) | '<inline json>'")
		return 2
	}

	body, err := readBody(args[2:])
	if err != nil {
		return fail(err)
	}

	res, err := api.Do(args[0], args[1], key, body)
	if err != nil {
		return fail(err)
	}

	fmt.Println(strings.TrimRight(string(res.Body), "\n"))
	if res.OK() {
		return 0
	}
	return 1
}

// readBody resolves an optional body argument into validated JSON bytes and
// fails fast if it isn't JSON (the API speaks JSON only).
func readBody(rest []string) ([]byte, error) {
	if len(rest) == 0 {
		return nil, nil
	}
	raw, source, err := loadBody(rest[0])
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("request body is not valid JSON (%s)", source)
	}
	return raw, nil
}

// loadBody turns a body argument into raw bytes plus a source label for errors.
// Resolution order:
//
//	"-"             read from stdin
//	"@path"         read the file at path (explicit; a missing file errors)
//	an existing path read that file (bare path, so `POST … topic.json` works)
//	anything else   treat the argument as a literal JSON string
func loadBody(arg string) (data []byte, source string, err error) {
	switch {
	case arg == "-":
		data, err = io.ReadAll(os.Stdin)
		return data, "stdin", err
	case strings.HasPrefix(arg, "@"):
		path := arg[1:]
		data, err = os.ReadFile(path)
		return data, "file " + path, err
	default:
		if info, statErr := os.Stat(arg); statErr == nil && !info.IsDir() {
			data, err = os.ReadFile(arg)
			return data, "file " + arg, err
		}
		return []byte(arg), "inline", nil
	}
}
