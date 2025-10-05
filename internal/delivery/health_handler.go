package delivery

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/mgmacri/pool-maintenance-app/internal/version"
    "go.uber.org/zap"
)

// HealthCheckResponse represents the response body for the health check endpoint.
//
// Example:
//
//	{
//	  "status": "ok",
//	  "version": "1.0.0",
//	  "commit": "abc1234",
//	  "build_date": "2025-08-25T12:34:56Z"
//	}
type HealthCheckResponse struct {
	Status    string `json:"status" example:"ok"`
	Version   string `json:"version" example:"1.0.0"`
	Commit    string `json:"commit" example:"abc1234"`
	BuildDate string `json:"build_date" example:"2025-08-25T12:34:56Z"`
}

// ReadinessResponse will evolve in later commits to include dependency checks; for now matches HealthCheckResponse.
type ReadinessResponse struct {
	Status    string `json:"status" example:"ok"`
	Version   string `json:"version" example:"1.0.0"`
	Commit    string `json:"commit" example:"abc1234"`
	BuildDate string `json:"build_date" example:"2025-08-25T12:34:56Z"`
}

// HealthHandler defines a handler for health checks.
type HealthHandler struct {
	Logger   *zap.Logger
	checkers []ReadinessChecker
}

// ReadinessChecker defines a dependency readiness contract. Future implementations might check database, cache, message broker, etc.
type ReadinessChecker interface {
	Name() string
	Check() error
}

// NewHealthHandler creates a new HealthHandler with the provided logger.
func NewHealthHandler(logger *zap.Logger, checkers ...ReadinessChecker) *HealthHandler {
	return &HealthHandler{Logger: logger, checkers: checkers}
}

// Check returns a simple health status and logs the health check.
//
// @Summary      Health check
// @Description  Returns service health and version info. Useful for uptime monitoring, CI/CD, and debugging.
// @Description
// @Description  **Example GitHub Actions usage:**
// @Description  A step in your CI/CD pipeline to verify deployment.
// @Description  ```yaml
// @Description  - name: Check service health
// @Description    uses: jtalk/url-health-check-action@v4
// @Description    with:
// @Description      url: https://your-app.com/health/live
// @Description      max-attempts: 10
// @Description      retry-delay: 5s
// @Description  ```
// @Tags         health
// @Produce      json
// @Success      200  {object}  delivery.HealthCheckResponse
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) { // legacy alias for backward compatibility
	h.Logger.Info("health check endpoint called (alias for /health/live)", zap.String("path", c.FullPath()))
	h.livePayload(c)
}

// Live returns liveness status (process is up). Designed to stay fast & allocation-light.
// @Summary Liveness probe
// @Description Returns 200 if the process is running. Avoids external dependency checks.
// @Tags health
// @Produce json
// @Success 200 {object} delivery.HealthCheckResponse
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	h.Logger.Debug("liveness probe", zap.String("path", c.FullPath()))
	h.livePayload(c)
}

// Ready returns readiness status. In this initial commit it mirrors liveness; later commits will add dependency evaluation.
// @Summary Readiness probe
// @Description Indicates whether the service is ready to accept traffic. Will include dependency statuses in later iterations.
// @Tags health
// @Produce json
// @Success 200 {object} delivery.ReadinessResponse
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	h.Logger.Debug("readiness probe", zap.String("path", c.FullPath()))
	// Evaluate placeholder checkers. If any fail, return 503 (degraded).
	for _, chk := range h.checkers {
		if err := chk.Check(); err != nil {
			h.Logger.Warn("readiness dependency failed", zap.String("dependency", chk.Name()), zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, ReadinessResponse{
				Status:    "degraded",
				Version:   version.Version,
				Commit:    version.Commit,
				BuildDate: version.BuildDate,
			})
			return
		}
	}
	c.JSON(http.StatusOK, ReadinessResponse{
		Status:    "ok",
		Version:   version.Version,
		Commit:    version.Commit,
		BuildDate: version.BuildDate,
	})
}

func (h *HealthHandler) livePayload(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:    "ok",
		Version:   version.Version,
		Commit:    version.Commit,
		BuildDate: version.BuildDate,
	})
}
