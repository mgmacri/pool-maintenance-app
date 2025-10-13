# Health Endpoint Contract

This document defines the health check endpoints for the Pool Maintenance API, following Kubernetes liveness and readiness probe conventions.

## Endpoints Overview

| Endpoint | Purpose | Status Codes | Use Case |
|----------|---------|--------------|----------|
| `/health` | Legacy alias for `/health/live` | 200 | Backward compatibility |
| `/health/live` | Liveness probe | 200 | Process is running |
| `/health/ready` | Readiness probe | 200, 503 | Service ready to accept traffic |

## Liveness Probe: `/health/live`

**Purpose:** Indicates whether the application process is alive and responsive.

**Success Criteria:** Process can handle HTTP requests.

**Response Format:**
```json
{
  "status": "ok",
  "version": "dev",
  "commit": "abc1234",
  "build_date": "2025-10-13T10:30:00Z"
}
```

**HTTP Status:** Always `200 OK` unless the process is completely unresponsive.

**Kubernetes Configuration:**
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3
```

## Readiness Probe: `/health/ready`

**Purpose:** Indicates whether the service is ready to handle business traffic by checking external dependencies.

**Success Criteria:** All registered dependency checkers pass.

**Response Formats:**

### Healthy (200 OK)
```json
{
  "status": "ok",
  "version": "dev", 
  "commit": "abc1234",
  "build_date": "2025-10-13T10:30:00Z",
  "dependencies": [
    {
      "name": "database",
      "status": "ok"
    },
    {
      "name": "cache",
      "status": "ok"
    }
  ]
}
```

### Degraded (503 Service Unavailable)
```json
{
  "status": "degraded",
  "version": "dev",
  "commit": "abc1234", 
  "build_date": "2025-10-13T10:30:00Z",
  "dependencies": [
    {
      "name": "database",
      "status": "ok"
    },
    {
      "name": "cache", 
      "status": "degraded",
      "error": "connection timeout after 5s"
    }
  ]
}
```

**HTTP Status:**
- `200 OK`: All dependencies healthy
- `503 Service Unavailable`: One or more dependencies failing

**Kubernetes Configuration:**
```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  successThreshold: 1
  failureThreshold: 2
```

## Dependency Status Values

| Status | Meaning |
|--------|---------|
| `ok` | Dependency is healthy and responding |
| `degraded` | Dependency is failing or unreachable |

## Implementation Notes

### Adding New Dependencies

Dependencies are registered via the `ReadinessChecker` interface:

```go
type ReadinessChecker interface {
    Name() string
    Check() error
}
```

Example implementation:
```go
type DatabaseChecker struct {
    db *sql.DB
}

func (d *DatabaseChecker) Name() string {
    return "database"
}

func (d *DatabaseChecker) Check() error {
    return d.db.Ping()
}
```

Register during service initialization:
```go
dbChecker := &DatabaseChecker{db: db}
healthHandler := delivery.NewHealthHandler(logger, dbChecker)
```

### Performance Considerations

- Liveness checks should be fast (< 100ms) and avoid external calls
- Readiness checks should timeout quickly (< 3s per dependency)  
- Consider circuit breakers for flaky dependencies
- Limit the number of dependencies to avoid cascading failures

### Monitoring Integration

Health endpoints provide structured data for monitoring:

- **Uptime monitoring:** Use `/health/live` for basic availability
- **Load balancer health:** Use `/health/ready` for traffic routing decisions
- **Dependency dashboards:** Parse `dependencies` array for per-service metrics
- **Alerting:** 503 responses indicate degraded state requiring investigation

## Legacy Compatibility

The `/health` endpoint serves as an alias to `/health/live` for backward compatibility. New integrations should prefer the explicit `/health/live` and `/health/ready` endpoints.

Future deprecation of `/health` will be communicated via:
1. Response headers (`X-Deprecated: true`)
2. Structured logging warnings
3. Documentation updates

## Security Considerations

Health endpoints are **publicly accessible** without authentication to support:
- Load balancer health checks
- Kubernetes probes  
- External monitoring systems

Sensitive dependency details (connection strings, internal IPs) are never exposed in health responses. Only dependency names and high-level status are returned.