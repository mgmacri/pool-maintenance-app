package delivery

import (
	"net/http"
	"time"

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
	Status        string  `json:"status" example:"ok"`
	Version       string  `json:"version" example:"1.0.0"`
	Commit        string  `json:"commit" example:"abc1234"`
	BuildDate     string  `json:"build_date" example:"2025-08-25T12:34:56Z"`
	UptimeSeconds float64 `json:"uptime_seconds" example:"123.45"`
}

// HealthHandler defines a handler for health checks.
type HealthHandler struct {
	Logger    *zap.Logger
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler with the provided logger.
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
    return &HealthHandler{Logger: logger, startTime: time.Now()}
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
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := time.Since(h.startTime).Seconds()
	h.Logger.Info("health check endpoint called", zap.String("path", c.FullPath()), zap.Float64("uptime_seconds", uptime))
	resp := HealthCheckResponse{
		Status:        "ok",
		Version:       version.Version,
		Commit:        version.Commit,
		BuildDate:     version.BuildDate,
		UptimeSeconds: uptime,
	}
	c.JSON(http.StatusOK, resp)
}
