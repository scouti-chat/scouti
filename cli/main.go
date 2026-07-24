// Command founderping is a tiny two-command CLI for the FounderPing devkit: it
// holds the developer's access key and forwards authenticated requests to the
// /api/v1 gate. All product logic lives server-side; this binary is deliberately thin.
package main

import (
	"fmt"
	"os"

	"github.com/founderping/founderping/cli/internal/commands"
)

// version is overwritten at build time via -ldflags "-X main.version=...".
var version = "dev"

const usage = `founderping — auth + request forwarder for FounderPing (see ../skill/SKILL.md)

Usage:
  founderping login [--token <UAK>]           Authorize this machine via the browser
                                              device flow, or store a pre-issued key
  founderping request <METHOD> <PATH> [body]  Call the /api/v1 gate with your key
                                                body: file.json | @file.json |
                                                      - (stdin) | '<inline json>'

Options:
  -h, --help       Show this help
  -v, --version    Show the CLI version

Environment:
  FOUNDERPING_API_URL      Override the API origin (staging / local dev)
  FOUNDERPING_ACCESS_KEY   Use this key instead of ~/.founderping/credentials.json (CI)
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Print(usage)
		return 0
	}
	switch args[0] {
	case "login":
		return commands.Login(args[1:])
	case "request":
		return commands.Request(args[1:])
	case "-v", "--version", "version":
		fmt.Println(version)
		return 0
	case "-h", "--help", "help":
		fmt.Print(usage)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
