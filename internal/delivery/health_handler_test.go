package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()
	h := NewHealthHandler(logger)
	r.GET("/health", h.Check)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse JSON response using encoding/json
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Assert required fields
	assert.Equal(t, "ok", resp["status"])
	assert.Contains(t, resp, "version")
	assert.Contains(t, resp, "commit")
	assert.Contains(t, resp, "build_date")
	assert.Contains(t, resp, "uptime_seconds")

	// Type assertions
	assert.IsType(t, "", resp["version"])
	assert.IsType(t, "", resp["commit"])
	assert.IsType(t, "", resp["build_date"])
	// JSON numbers unmarshal to float64
	assert.IsType(t, float64(0), resp["uptime_seconds"])
	// Sanity: uptime should be >= 0
	assert.GreaterOrEqual(t, resp["uptime_seconds"].(float64), 0.0)
}
