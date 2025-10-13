package version

// Package version provides build-time metadata injected via Go linker flags.
//
// Expected ldflags (example):
//   -X 'github.com/mgmacri/pool-maintenance-app/internal/version.Version=0.1.0'
//   -X 'github.com/mgmacri/pool-maintenance-app/internal/version.Commit=abc1234'
//   -X 'github.com/mgmacri/pool-maintenance-app/internal/version.BuildDate=2025-10-05T12:34:56Z'
//
// Defaults below apply when built locally without those flags (e.g., `go run`).
// CI / Docker build overrides them for reproducible artifacts.
var (
	// Version is the semantic version of the binary (e.g. 0.1.0). Defaults to "dev" when not set.
	Version = "dev"
	// Commit is the git SHA of the source used to build the binary.
	Commit = ""
	// BuildDate is the RFC3339 timestamp of the build.
	BuildDate = ""
)
