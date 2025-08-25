package version

// These variables are set at build time using -ldflags.
var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)
