// Package version implements a customizable version for the MOAC binaries.
package version

import (
	"runtime/debug"
)

// version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var version string

// GetVersion fetches the version of the MOAC binaries, configurable at link-time.
func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && version == "" {
		return info.Main.Version
	}

	return version
}
