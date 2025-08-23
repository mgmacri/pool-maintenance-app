# pool-maintenance-app

Initial Go project setup for the Pool Maintenance API.

## Getting Started

1. Ensure you have Go 1.25+ installed.
2. Clone the repository.
3. Install dependencies:
	```sh
	go mod tidy
	```
4. Run the application:
	```sh
	go run ./cmd/main.go
	```

## Project Structure

- `cmd/` - Application entry point
- `internal/` - Clean architecture layers (domain, usecase, delivery, repository)
- `pkg/` - Shared utilities
- `docs/` - Documentation

---
This is the initial setup commit. Further features will be added incrementally.
