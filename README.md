# pool-maintenance-app

![Go CI](https://github.com/mgmacri/pool-maintenance-app/actions/workflows/go-ci.yml/badge.svg)

Pool Maintenance API — a Go project following Clean Architecture and DevOps best practices.



## API Documentation (Swagger/OpenAPI)

This project uses [Swagger/OpenAPI](https://swagger.io/) for interactive API documentation, generated with [swaggo/swag](https://github.com/swaggo/swag).
# pool-maintenance-app

![Go CI](https://github.com/mgmacri/pool-maintenance-app/actions/workflows/go-ci.yml/badge.svg)

> Open-source **portfolio & learning project**: a narrative codebase showing how to evolve a vertical slice (health endpoint → observability → domain) with Clean Architecture, DevOps/SRE, and documentation discipline. **Not production-certified.** See `DISCLAIMER.md`.

## Purpose & Vision
- Incremental observability-first development
- Clean layering (delivery / use case / domain)
- CI/CD, security scanning, SBOM & signing (planned)
- Evolving requirements: CRS (`design/crs.md`) ↔ ERS (`design/ers.md`)
- Chapter-based delivery plan (`plan.md`)

## Roadmap Snapshot (See `plan.md`)
| Chapter | Theme | Focus Value |
|---------|-------|-------------|
| 1 | Observability Slice | Health endpoints, metrics, tracing scaffold |
| 2 | Security & RBAC | Auth, roles, audit trail |
| 3 | Chemistry Core | Test input → dose engine interface → report |
| 4 | Scheduling & Routing | Service plans, manual route ordering |
| 5 | Billing & Exports | Invoicing + compliance export |

## API Documentation (Swagger / OpenAPI)
Generated via [swaggo/swag](https://github.com/swaggo/swag).

View locally (when server running): http://localhost:8080/swagger/index.html

Regenerate after changing annotated comments:
```sh
swag init -g cmd/main.go
```
OpenAPI artifacts live in `docs/` and are bundled in the Docker image.

## Run Locally (Go)
```sh
go mod tidy
go run ./cmd/main.go
```

## Build Metadata (Version, Commit, Build Date, Uptime)
The binary embeds build-time metadata surfaced at `/health`:

| Field | Source | Purpose |
|-------|--------|---------|
| `version` | `-ldflags` (or defaults to `dev`) | Human + automation friendly release identifier |
| `commit` | Git short SHA | Precise reproducibility / traceability |
| `build_date` | UTC RFC3339 timestamp | Audit / release notes correlation |
| `uptime_seconds` | In-process runtime | Liveness diagnostics, quick sanity |

### Local Development
Running with `go run` (or plain `go build`) will show:
```json
{"version":"dev","commit":"","build_date":""}
```
CI / Docker builds inject real values.

### Makefile
```sh
make build            # embeds version, commit, build date
make run              # builds then runs
make info             # prints the resolved ldflags values
```
Override version explicitly:
```sh
make build VERSION=0.1.0
```

### Manual go build Example
```sh
go build -ldflags "-X 'github.com/mgmacri/pool-maintenance-app/internal/version.Version=0.1.0' \
	-X 'github.com/mgmacri/pool-maintenance-app/internal/version.Commit=$(git rev-parse --short HEAD)' \
	-X 'github.com/mgmacri/pool-maintenance-app/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
	-o bin/pool-maintenance-api ./cmd/main.go
```

### Sample Health Endpoint Response
```bash
curl -s http://localhost:8080/health | jq
```
Possible output (locally without ldflags):
```json
{
	"status": "ok",
	"version": "dev",
	"commit": "",
	"build_date": "",
	"uptime_seconds": 0.123456
}
```
In CI-built image these fields will be populated with release values.

## Run with Docker (Static Alpine Build)
```sh
docker build -t pool-maintenance-api .
docker run -p 8080:8080 pool-maintenance-api
```
Health endpoint: http://localhost:8080/health

> Image: statically linked (musl) for portability.

## Observability & Logging (Current State)
Structured logging with [zap](https://github.com/uber-go/zap). Fields include `service`, `env`, `version`, and a placeholder `trace_id`. Tracing & richer metrics will be added per Chapter 1 tasks.

Example log (abbreviated):
```json
{"level":"info","msg":"request completed","path":"/health","trace_id":"","version":"dev"}
```

## Project Structure
| Path | Purpose |
|------|---------|
| `cmd/` | Application entrypoint |
| `internal/delivery/` | HTTP handlers (REST) |
| `internal/usecase/` | Business orchestration layer (future expansion) |
| `internal/domain/` | Core entities & domain logic stubs |
| `internal/middleware/` | Cross-cutting HTTP middleware (logging, future tracing) |
| `docs/` | Generated Swagger + doc assets |
| `design/` | CRS / ERS specifications |
| `plan.md` | Iterative delivery & blog plan |
| `DISCLAIMER.md` | Portfolio / non-production notice |

## Contributing (Learning-Focused)
This is an educational repository. PRs are welcome when they:
- Improve clarity (docs, structure, tests) OR
- Advance a planned chapter task from `plan.md`.

Workflow:
1. Branch: `feat/<short-desc>`
2. Ensure `go test ./...` passes
3. Keep commits small & conventional (e.g., `feat:`, `docs:`)
4. Reference relevant ERS IDs in PR description if implementing requirements.

## Running Tests
```sh
go test ./...
```

## (Planned) Local Actions / CI Helpers
```sh
act -j build   # Run GitHub Actions locally (optional)
```

## License & Disclaimer
Licensed under [MIT](LICENSE). See `DISCLAIMER.md` for limitations and intended use.

## Blog / Learning Series (Planned)
| # | Working Title | Status |
|---|---------------|--------|
| 1 | From Zero to Production-Ready Health Endpoint | Drafting |
| 2 | Securing the Backbone: Auth & RBAC Foundations | Pending |
| 3 | Chemistry Intelligence: Designing a Dose Engine Interface | Pending |
| 4 | Scheduling & Route Foundations | Pending |
| 5 | Monetizing via Reports & Billing (Demonstration) | Pending |

---
**Questions / Ideas?** Open an issue explaining the learning value.
Note: The artifact upload step is skipped locally, and Trivy or golangci-lint must be installed in the runner image. Security scanning may fail the build if vulnerabilities are found—this is intentional for best practices.
