package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler defines a handler for health checks.
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check returns a simple health status.
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
