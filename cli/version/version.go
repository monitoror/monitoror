//+build !faker

package version

// Default build-time variable.
// These values are overridden via ldflags on build
var (
	Version   = "v0.0.0-dev"
	GitCommit = "unknown-gitcommit"
	BuildTime = "unknown-buildtime"
)
