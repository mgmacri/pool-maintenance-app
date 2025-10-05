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

func TestHealthHandler_Check_LegacyAlias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := zap.NewNop()
	h := NewHealthHandler(logger)
	r.GET("/health", h.Check)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	validateStandardFields(t, w.Body.Bytes())
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
