package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/scouti-chat/scouti/cli/internal/creds"
	"github.com/scouti-chat/scouti/cli/internal/device"
)

// Login authorizes this machine and stores the resulting access key.
//
// With --token (or SCOUTI_ACCESS_KEY) it stores a pre-issued key for CI /
// headless use. Otherwise it runs the device-authorization flow: open a browser
// to sign in and approve, then poll with the private device_code until the key
// is issued.
func Login(args []string) int {
	if token := tokenArg(args); token != "" {
		return save(creds.Credentials{AccessKey: token})
	}

	start, err := device.Begin()
	if err != nil {
		return fail(err)
	}

	fmt.Printf("\nOpen this URL to authorize the CLI:\n\n  %s\n\n", start.VerificationURIComplete)
	fmt.Printf("Confirm this code matches the page: %s\n\n", start.UserCode)
	if !device.OpenBrowser(start.VerificationURIComplete) {
		fmt.Println("(Could not open a browser automatically — copy the link above.)")
	}
	fmt.Println("Waiting for approval…")

	deadline := time.Now().Add(time.Duration(start.ExpiresIn) * time.Second)
	interval := start.Interval
	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)
		res, err := device.Poll(start.DeviceCode)
		if err != nil {
			return fail(err)
		}
		if res.Issued {
			return save(creds.Credentials{AccessKey: res.AccessKey, UserID: res.UserID})
		}
		interval = res.Interval
	}
	fmt.Fprintln(os.Stderr, "Login timed out before approval. Run `scouti login` again.")
	return 1
}

func save(c creds.Credentials) int {
	if err := creds.Write(c); err != nil {
		return fail(err)
	}
	fmt.Println("Logged in. Credentials saved to ~/.scouti/credentials.json")
	return 0
}

func tokenArg(args []string) string {
	for i, a := range args {
		if a == "--token" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return os.Getenv("SCOUTI_ACCESS_KEY")
}
