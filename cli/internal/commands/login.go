package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/scouti-chat/scouti/cli/internal/creds"
	"github.com/scouti-chat/scouti/cli/internal/device"
)

// maxLoginWait caps how long the device-authorization flow waits for the
// developer to approve in the browser. The server's code lives longer, but a
// forgotten or headless `scouti login` shouldn't hang the terminal — we give up
// and tell them to re-run it.
const maxLoginWait = 5 * time.Minute

// Login authorizes this machine and stores the resulting access key.
//
// With --token (or SCOUTI_ACCESS_KEY) it stores a pre-issued key for CI /
// headless use. Otherwise it runs the device-authorization flow: open a browser
// to sign in and approve, then poll with the private device_code until the key
// is issued or maxLoginWait elapses.
func Login(args []string) int {
	if token := tokenArg(args); token != "" {
		return save(creds.Credentials{AccessKey: token})
	}

	start, err := device.Begin(context.Background())
	if err != nil {
		return fail(err)
	}

	fmt.Printf("\nOpen this URL to authorize the CLI:\n\n  %s\n\n", start.VerificationURIComplete)
	fmt.Printf("Confirm this code matches the page: %s\n\n", start.UserCode)
	if !device.OpenBrowser(start.VerificationURIComplete) {
		fmt.Println("(Could not open a browser automatically — copy the link above.)")
	}
	fmt.Println("Waiting for approval…")

	// Bound the wait: never longer than maxLoginWait, nor past the server's expiry.
	wait := maxLoginWait
	if exp := time.Duration(start.ExpiresIn) * time.Second; exp > 0 && exp < wait {
		wait = exp
	}
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	interval := start.Interval
	for {
		select {
		case <-ctx.Done():
			return timedOut(wait)
		case <-time.After(time.Duration(interval) * time.Second):
		}
		res, err := device.Poll(ctx, start.DeviceCode)
		if err != nil {
			if ctx.Err() != nil {
				return timedOut(wait)
			}
			return fail(err)
		}
		if res.Issued {
			return save(creds.Credentials{AccessKey: res.AccessKey, UserID: res.UserID})
		}
		interval = res.Interval
	}
}

// timedOut reports that approval didn't arrive within wait and returns exit 1.
func timedOut(wait time.Duration) int {
	fmt.Fprintf(os.Stderr, "Login timed out after %s without approval. Run `scouti login` again.\n", wait)
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
