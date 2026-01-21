package version

// Version information injected at build time via ldflags
var (
	// Tag is the git tag or version string
	Tag = "dev"
	// Commit is the git commit SHA
	Commit = "unknown"
	// BuildTime is the build timestamp
	BuildTime = "unknown"
)
