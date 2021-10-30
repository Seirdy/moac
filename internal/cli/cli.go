// Package cli contains functions shared between moac and moac-pwgen binaries
package cli

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/rivo/uniseg"
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

var (
	// version can be set at link time to override debug.BuildInfo.Main.Version,
	// which is "(devel)" when building from within the module. See
	// golang.org/issue/29814 and golang.org/issue/29228.
	version = "(devel)"

	// ErrBadCmdline indicates an invalid argument has been passed via the CLI.
	ErrBadCmdline = errors.New("bad arguments")
)

// GetVersion fetches the version of the MOAC binaries, configurable at link-time.
func GetVersion() string {
	versionUnset := version == "" || version == "(devel)"
	if info, ok := debug.ReadBuildInfo(); ok && versionUnset {
		return info.Main.Version
	}

	return version
}

// HasGrapheme returns true if a string contains any grapheme clusters, false otherwise.
func HasGrapheme(str string) bool {
	graphemes := uniseg.NewGraphemes(str)
	for graphemes.Next() {
		if len(graphemes.Runes()) > 1 {
			return true
		}
	}

	return false
}
