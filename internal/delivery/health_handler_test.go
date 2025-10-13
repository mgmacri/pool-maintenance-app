package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthHandler_Check_LegacyAlias(t *testing.T) {
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

func TestHealthHandler_Live(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()
	h := NewHealthHandler(logger)
	r.GET("/health/live", h.Live)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/live", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	validateStandardFields(t, w.Body.Bytes())
}

func TestHealthHandler_Ready_InitialMirrorsLive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()
	h := NewHealthHandler(logger)
	r.GET("/health/ready", h.Ready)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/ready", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	validateStandardFields(t, w.Body.Bytes())
}

func TestHealthHandler_Ready_Returns503OnDependencyFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()

	// Create handler with failing checker
	failingChecker := &fakeChecker{name: "db", shouldFail: true}
	h := NewHealthHandler(logger, failingChecker)
	r.GET("/health/ready", h.Ready)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/ready", nil)
	r.ServeHTTP(w, req)

	// Should return 503 Service Unavailable
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp ReadinessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "degraded", resp.Status)
	assert.Len(t, resp.Dependencies, 1)
	assert.Equal(t, "db", resp.Dependencies[0].Name)
	assert.Equal(t, "degraded", resp.Dependencies[0].Status)
	assert.Equal(t, "simulated failure", resp.Dependencies[0].Error)
}

func TestHealthHandler_Ready_MixedDependencyStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()

	// Mix of healthy and failing checkers
	healthyChecker := &fakeChecker{name: "cache", shouldFail: false}
	failingChecker := &fakeChecker{name: "db", shouldFail: true}
	h := NewHealthHandler(logger, healthyChecker, failingChecker)
	r.GET("/health/ready", h.Ready)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/ready", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp ReadinessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "degraded", resp.Status)
	assert.Len(t, resp.Dependencies, 2)

	// Check both dependencies are reported
	depsByName := make(map[string]DependencyStatus)
	for _, dep := range resp.Dependencies {
		depsByName[dep.Name] = dep
	}

	assert.Equal(t, "ok", depsByName["cache"].Status)
	assert.Empty(t, depsByName["cache"].Error)
	assert.Equal(t, "degraded", depsByName["db"].Status)
	assert.Equal(t, "simulated failure", depsByName["db"].Error)
}

// fakeChecker implements ReadinessChecker for testing
type fakeChecker struct {
	name       string
	shouldFail bool
}

func (f *fakeChecker) Name() string {
	return f.name
}

func (f *fakeChecker) Check() error {
	if f.shouldFail {
		return errors.New("simulated failure")
	}
	return nil
}

func validateStandardFields(t *testing.T, body []byte) {
	var resp map[string]interface{}
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Contains(t, resp, "version")
	assert.Contains(t, resp, "commit")
	assert.Contains(t, resp, "build_date")
	assert.IsType(t, "", resp["version"])
	assert.IsType(t, "", resp["commit"])
	assert.IsType(t, "", resp["build_date"])
}
