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

This project uses [Uber Zap](https://github.com/uber-go/zap) for structured, production-grade JSON logging. All logs include:

- `service`: the service name (e.g., `pool-maintenance-api`)
- `env`: the environment (from the `ENV` environment variable, defaults to `dev`)
- `version`: the build version (from ldflags)
- `trace_id`: a placeholder for distributed tracing (currently null, will be populated when tracing is integrated)

**Log Example:**
```json
{
	"level": "info",
	"ts": 1692979200.123,
	"caller": "internal/middleware/zap.go:20",
	"msg": "request completed",
	"service": "pool-maintenance-api",
	"env": "dev",
	"version": "dev",
	"trace_id": null,
	"status": 200,
	"method": "GET",
	"path": "/health",
	"ip": "127.0.0.1",
	...
}
```

These fields make logs easy to aggregate and search in systems like Loki, Elasticsearch, or Datadog. Similarly, the `trace_id` field will integrate with tracing in a later feature update. 

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
