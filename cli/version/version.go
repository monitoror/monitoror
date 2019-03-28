//+build !faker

package version

// Default build-time variable.
// These values are overridden via ldflags on build
var (
	Version   = "x.x.x-dev"
	GitCommit = "unknown-gitcommit"
	BuildTime = "unknown-buildtime"
)
