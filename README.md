# pool-maintenance-app

![Go CI](https://github.com/mgmacri/pool-maintenance-app/actions/workflows/go-ci.yml/badge.svg)

Pool Maintenance API — a Go project following Clean Architecture and DevOps best practices.



## API Documentation (Swagger/OpenAPI)

This project uses [Swagger/OpenAPI](https://swagger.io/) for interactive API documentation, generated with [swaggo/swag](https://github.com/swaggo/swag).

- **View the docs:**
	- Locally: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
	- In Docker: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Regenerate docs after changing endpoint comments:**
	```sh
	swag init -g cmd/main.go
	```
- **Swagger files are in the `docs/` directory and are copied into the Docker image.**

---

### Run Locally (Go)
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


### Run with Docker (Static musl/Alpine build)
1. Build the Docker image (now uses Alpine and a fully static musl-linked Go binary):
	```sh
	docker build -t pool-maintenance-api .
	```
2. Run the container:
	```sh
	docker run -p 8080:8080 pool-maintenance-api
	```
3. Access the health check endpoint:
	[http://localhost:8080/health](http://localhost:8080/health)

**Note:** The Docker image is now based on Alpine Linux and contains a statically linked Go binary built with musl libc. This eliminates glibc version issues (e.g., "GLIBC_x.x not found") and ensures maximum portability across Linux hosts. See the Dockerfile for build details.




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

---
## Project Structure

- `cmd/` — Application entry point
- `internal/delivery/` — HTTP handlers (e.g., health check handler)
- `internal/` — Clean architecture layers (domain, usecase, repository)
- `pkg/` — Shared utilities
- `docs/` — Documentation


## Contributing

We follow an industry-standard Git workflow:

1. Create a new branch for each feature or fix (e.g., `feat/feature-name`, `ci/add-go-test-step`).
2. Make your changes and commit with clear, conventional messages.
3. Push your branch and open a Pull Request (PR) to `main`.
4. All PRs require at least one review and must pass CI checks before merging.
5. After merging, delete the feature branch if no longer needed.

## Running Tests

To run all tests locally:
```sh
go test ./...
```



```sh
act -j build
```

Note: The artifact upload step is skipped locally, and Trivy or golangci-lint must be installed in the runner image. Security scanning may fail the build if vulnerabilities are found—this is intentional for best practices.



## License

This project is licensed under the terms of the [MIT License](LICENSE).

---
For questions or contributions, open an issue or PR.
