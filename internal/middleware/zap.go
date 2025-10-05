package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
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
	traceID := extractTraceID(c)
	c.Set("trace_id", traceID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
	// traceID already extracted (empty if none provided)
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
const (
	traceParentHeader = "traceparent"     // W3C
	b3TraceIDHeader   = "X-B3-TraceId"    // B3 single header variant (only trace id here)
)

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

var (
	reTraceParent = regexp.MustCompile(`^[ \t]*([0-9a-f]{2})-([0-9a-f]{32})-([0-9a-f]{16})-([0-9a-f]{2})(?:-.*)?$`)
	reHex16Or32   = regexp.MustCompile(`^[0-9a-f]{16}|[0-9a-f]{32}$`)
)

// extractTraceID attempts to pull a trace id from standard propagation headers.
// Precedence: W3C traceparent (if valid) > B3 TraceId > empty string.
func extractTraceID(c *gin.Context) string {
	if tp := c.GetHeader(traceParentHeader); tp != "" {
		if id := parseTraceParent(tp); id != "" { return id }
	}
	if b3 := c.GetHeader(b3TraceIDHeader); b3 != "" {
		b3 = strings.TrimSpace(b3)
		if reHex16Or32.MatchString(b3) { return b3 }
	}
	return ""
}

// parseTraceParent validates and returns the 32-hex trace-id from a W3C traceparent header.
func parseTraceParent(headerVal string) string {
	m := reTraceParent.FindStringSubmatch(strings.TrimSpace(headerVal))
	if len(m) != 5 { return "" }
	version, traceID, parentID, flags := m[1], m[2], m[3], m[4]
	// Basic sanity: ensure not all zeros (per spec recommendation)
	if allZeros(traceID) || allZeros(parentID) { return "" }
	_ = version; _ = flags // not used now, but parse retained for future logic
	return traceID
}

func allZeros(s string) bool {
	for _, r := range s { if r != '0' { return false } }
	return true
}
