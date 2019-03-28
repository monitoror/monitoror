//+build faker

package version

// Default build-time variable.
// Fixed for faker server
var (
	Version   = "x.x.x-faker"
	GitCommit = "commit-hash"
	BuildTime = "2017-01-09 22:45:00+00:00"
)
