// Package cli contains functions shared between moac and moac-pwgen binaries
package cli

import (
	"fmt"
	"os"
)

// FloatFmt defines how many digits of a float to print.
const FloatFmt = "%.3g\n"

// ExitOnErr exits the program with status 1 with a message in the presence of an error.
func ExitOnErr(err error, extraLine string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %s\n%s", err.Error(), extraLine)
		os.Exit(1)
	}
}
