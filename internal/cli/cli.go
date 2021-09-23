// Package cli contains functions shared between moac and moac-pwgen binaries
package cli

import (
	"fmt"
	"os"
)

// ExitOnErr exits the program with status 1 with a message in the presence of an error.
func ExitOnErr(err error, extraLine string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %s\n%s", err.Error(), extraLine)
		os.Exit(1)
	}
}
