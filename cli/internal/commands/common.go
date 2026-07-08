// Package commands implements the scouti subcommands (login, request).
package commands

import (
	"fmt"
	"os"
)

// fail reports err on stderr and returns the process exit code for failures.
func fail(err error) int {
	fmt.Fprintln(os.Stderr, "Error:", err)
	return 1
}
