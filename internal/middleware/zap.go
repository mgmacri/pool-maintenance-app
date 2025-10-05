package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger returns a Gin middleware that logs requests using the provided zap.Logger.
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Correlation: request_id from incoming header or generated.
		reqID := getOrCreateRequestID(c)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		// Placeholder for trace_id, to be populated when tracing is integrated
		traceID := ""
		logger.Info("request completed",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("request_id", reqID),
			zap.String("trace_id", traceID),
		)
	}
}

const requestIDHeader = "X-Request-ID"

// getOrCreateRequestID returns the incoming request id or generates a new one.
// The chosen id is set on the response header and stored in the Gin context.
func getOrCreateRequestID(c *gin.Context) string {
	rid := c.GetHeader(requestIDHeader)
	if rid == "" {
		rid = newRequestID()
	}
	c.Header(requestIDHeader, rid)
	c.Set("request_id", rid)
	return rid
}

// newRequestID generates a 16-byte random hex string (32 chars) for correlation.
func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-derived value if RNG fails (extremely unlikely)
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(b)
}
