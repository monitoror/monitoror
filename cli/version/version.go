package version

// Default build-time variable.
// These values are overridden via ldflags on build
var (
	Version   = "unknown-version"
	GitCommit = "unknown-gitcommit"
	BuildTime = "unknown-buildtime"
)
