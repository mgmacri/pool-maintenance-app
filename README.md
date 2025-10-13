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

## Health Endpoints

The service provides Kubernetes-compatible health endpoints:

- `/health` - Legacy endpoint (alias for `/health/live`)
- `/health/live` - Liveness probe (fast, no external dependencies)
- `/health/ready` - Readiness probe (includes dependency checks)

See [docs/health.md](docs/health.md) for detailed API documentation and Kubernetes integration examples.

## Build Metadata (Version, Commit, Build Date, Uptime)
The binary embeds build-time metadata surfaced at health endpoints:

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

## Run with Docker (Static Alpine Build)
```sh
docker build -t pool-maintenance-api .
docker run -p 8080:8080 pool-maintenance-api
```
Health endpoints: 
- http://localhost:8080/health
- http://localhost:8080/health/live
- http://localhost:8080/health/ready

> Image: statically linked (musl) for portability.

## Observability & Logging

This project uses [Uber Zap](https://github.com/uber-go/zap) for structured JSON logging with early correlation primitives in place. Each request is logged once on completion with standardized fields to enable aggregation, filtering, and future distributed tracing.

### Correlation Fields

Included in each request log:

| Field | Source | Description |
|-------|--------|-------------|
| `service` | build-time constant | Logical service identifier |
| `env` | `ENV` env var (default `dev`) | Deployment environment tag |
| `version` | ldflags (`-X internal/version.Version`) | Git or semantic build version |
| `request_id` | Incoming `X-Request-ID` header or generated | Stable per-request correlation id (32 hex chars if generated) |
| `trace_id` | Incoming W3C `traceparent` or `X-B3-TraceId` | Distributed trace identifier (empty if not provided) |
| `status`, `method`, `path`, `latency` | HTTP layer | Request outcome metadata |

### Request ID Behavior
If a client sends `X-Request-ID`, it is preserved and echoed back. If absent, a 16 byte (32 hex) cryptographically random identifier is generated and returned in the same header. Always propagate this header across downstream calls inside scripts, batch jobs, or tests to stitch cross-service logs.

Example (no incoming request id):
```bash
curl -i http://localhost:8080/health | grep -i x-request-id
```

### Trace ID Extraction (Early Tracing Readiness)
The middleware inspects (precedence order):
1. `traceparent` (W3C Trace Context) – validates format `00-<32 hex trace id>-<16 hex parent id>-<2 hex flags>` and ignores invalid / all-zero fields.
2. `X-B3-TraceId` (B3 single field) – accepts 16 or 32 lowercase hex.

If present & valid, the `trace_id` field is added to the structured log. Otherwise the field is an empty string. Full tracing (spans/export) will come in a later instrumentation sub-chapter.

Example using `traceparent`:
```bash
curl -H "traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01" \
     http://localhost:8080/health
```

### Log Levels
Runtime log verbosity is controlled by the `LOG_LEVEL` environment variable (case-insensitive). Supported values:
`debug`, `info` (default), `warn` / `warning`, `error` / `err`, `dpanic`, `panic`, `fatal`.

Invalid values fall back silently to `info` (future enhancement: emit a startup warning). Example:
```bash
LOG_LEVEL=debug go run ./cmd/main.go
```

Verify debug suppression vs emission (excerpt using tests / observer core):
```bash
go test -run TestZapLogger_DebugLogFiltering ./internal/middleware -v
```

### Sample Log (Redacted for Brevity)
```json
{
  "level": "info",
  "msg": "request completed",
  "service": "pool-maintenance-api",
  "env": "dev",
  "version": "dev",
  "request_id": "f3b1a6d9099f4f42e8e97d5d6d3fe0c2",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "status": 200,
  "method": "GET",
  "path": "/health",
  "latency": 0.0012345
}
```

### Future Enhancements (Planned)
- Inject OpenTelemetry SDK & exporter
- Span context propagation & sampling controls
- Log -> Trace correlation enrichment (`trace_flags`, `span_id`)
- Structured error classification & error code taxonomy

These foundations allow immediate value (correlation, filtering) while minimizing future refactor risk.

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
