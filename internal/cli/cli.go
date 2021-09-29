// Package cli contains functions shared between moac and moac-pwgen binaries
package cli

import (
	"fmt"
	"os"
)

// FloatFmt defines how many digits of a float to print.
const FloatFmt = "%.3g\n"

// DisplayErr prints an error to stdout and returns true if it's nil.
func DisplayErr(err error, extraLine string) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %s\n%s", err.Error(), extraLine)

		return false
	}

	return true
}
