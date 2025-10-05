package middleware

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "go.uber.org/zap/zaptest/observer"
    "github.com/stretchr/testify/assert"
)

// helper to build a test logger capturing entries
func testLogger() (*zap.Logger, *observer.ObservedLogs) {
    core, logs := observer.New(zapcore.InfoLevel)
    return zap.New(core), logs
}

func TestZapLogger_GeneratesRequestID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    logger, logs := testLogger()
    r := gin.New()
    r.Use(ZapLogger(logger))
    r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    rid := w.Header().Get("X-Request-ID")
    assert.NotEmpty(t, rid, "expected generated request id header")

    entries := logs.All()
    assert.NotEmpty(t, entries)
    found := false
    for _, e := range entries {
        if e.Message == "request completed" {
            // marshal fields to map for easier lookup
            b, _ := json.Marshal(e.ContextMap())
            var m map[string]interface{}
            _ = json.Unmarshal(b, &m)
            v, ok := m["request_id"]
            assert.True(t, ok, "log should contain request_id")
            assert.Equal(t, rid, v)
            found = true
        }
    }
    assert.True(t, found, "expected to find request completed log entry")
}

func TestZapLogger_RespectsIncomingRequestID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    logger, logs := testLogger()
    r := gin.New()
    r.Use(ZapLogger(logger))
    r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    req.Header.Set("X-Request-ID", "fixed-id-123")
    r.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "fixed-id-123", w.Header().Get("X-Request-ID"))

    entries := logs.All()
    for _, e := range entries {
        if e.Message == "request completed" {
            b, _ := json.Marshal(e.ContextMap())
            var m map[string]interface{}
            _ = json.Unmarshal(b, &m)
            assert.Equal(t, "fixed-id-123", m["request_id"])
            return
        }
    }
    t.Fatalf("request completed log entry not found")
}
