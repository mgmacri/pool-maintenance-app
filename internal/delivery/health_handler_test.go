package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	// First request
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var resp1 map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.NoError(t, err)

	// Basic field assertions first call
	assert.Equal(t, "ok", resp1["status"])
	assert.IsType(t, "", resp1["version"])
	assert.IsType(t, "", resp1["commit"])
	assert.IsType(t, "", resp1["build_date"])
	assert.IsType(t, float64(0), resp1["uptime_seconds"])
	firstUptime := resp1["uptime_seconds"].(float64)
	assert.GreaterOrEqual(t, firstUptime, 0.0)

	// Wait a tiny bit to ensure uptime advances
	time.Sleep(25 * time.Millisecond)

	// Second request
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var resp2 map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.NoError(t, err)

	// Uptime should be monotonic increasing
	assert.Greater(t, resp2["uptime_seconds"].(float64), firstUptime)

	// Version default expectation when built without ldflags
	// Accept "dev" or any non-empty string (future builds). Just ensure key sanity.
	assert.NotNil(t, resp2["version"])
	// Build date may be empty if not injected; key should exist.
	assert.Contains(t, resp2, "build_date")
}
