package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

func TestZapLogger_ExtractsTraceID_FromTraceparent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, logs := testLogger()
	r := gin.New()
	r.Use(ZapLogger(logger))
	r.GET("/tp", func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	// traceparent: version(00)-traceid(32hex)-parentid(16hex)-flags(2hex)
	req, _ := http.NewRequest("GET", "/tp", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	entries := logs.All()
	for _, e := range entries {
		if e.Message == "request completed" {
			b, _ := json.Marshal(e.ContextMap())
			var m map[string]interface{}
			_ = json.Unmarshal(b, &m)
			assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", m["trace_id"])
			return
		}
	}
	t.Fatalf("trace_id not found in logs")
}

func TestZapLogger_ExtractsTraceID_FromB3(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, logs := testLogger()
	r := gin.New()
	r.Use(ZapLogger(logger))
	r.GET("/b3", func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/b3", nil)
	req.Header.Set("X-B3-TraceId", "463ac35c9f6413ad48485a3953bb6124")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	entries := logs.All()
	for _, e := range entries {
		if e.Message == "request completed" {
			b, _ := json.Marshal(e.ContextMap())
			var m map[string]interface{}
			_ = json.Unmarshal(b, &m)
			assert.Equal(t, "463ac35c9f6413ad48485a3953bb6124", m["trace_id"])
			return
		}
	}
	t.Fatalf("trace_id not found in logs for B3 header")
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

func TestZapLogger_DebugLogFiltering(t *testing.T) {
	gin.SetMode(gin.TestMode)

	run := func(level zapcore.Level) (hasDebug bool, hasInfo bool) {
		core, obs := observer.New(level)
		logger := zap.New(core)
		r := gin.New()
		r.Use(ZapLogger(logger))
		r.GET("/filter", func(c *gin.Context) {
			logger.Debug("debug diagnostic", zap.String("k", "v"))
			logger.Info("info diagnostic", zap.String("k", "v"))
			c.Status(204)
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/filter", nil)
		r.ServeHTTP(w, req)
		for _, e := range obs.All() {
			switch e.Message {
			case "debug diagnostic":
				hasDebug = true
			case "info diagnostic":
				hasInfo = true
			}
		}
		return
	}

	// Info level should NOT emit debug message
	dbg, info := run(zapcore.InfoLevel)
	if dbg {
		t.Fatalf("expected no debug log at info level")
	}
	if !info {
		t.Fatalf("expected info log at info level")
	}

	// Debug level should emit both
	dbg, info = run(zapcore.DebugLevel)
	if !dbg || !info {
		t.Fatalf("expected both debug and info logs at debug level (got debug=%v info=%v)", dbg, info)
	}
}
