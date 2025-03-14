package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDHeader     = "X-Trace-Id"
	TraceParentHeader = "traceparent"
	GinTraceIDKey     = "trace_id"
)

// FromContext retrieves trace ID from context
func FromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(GinTraceIDKey)
	if v == nil {
		return "", false
	}
	return v.(string), true
}

// generateTraceID generates a new valid trace ID
func generateTraceID() string {
	// Generate 16 random bytes (128 bits) for trace ID
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(randBytes)
}

// parseHexTraceID converts hex string to TraceID
func parseHexTraceID(hexStr string) (trace.TraceID, error) {
	if len(hexStr) != 32 {
		return trace.TraceID{}, fmt.Errorf("invalid trace id length")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return trace.TraceID{}, err
	}

	var tid trace.TraceID
	copy(tid[:], bytes)
	return tid, nil
}

// TraceID returns a middleware that handles trace ID propagation
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceID string
		ctx := c.Request.Context()
		spanContext := trace.SpanContextFromContext(ctx)

		// If there's a valid span context, use its trace ID
		if spanContext.IsValid() {
			traceID = spanContext.TraceID().String()
		} else if tid := c.GetHeader(TraceIDHeader); tid != "" {
			// Use trace ID from request header
			traceID = tid
		} else {
			// Generate new trace ID if none exists
			traceID = generateTraceID()
		}

		if traceID != "" {
			// Create new span context with the trace ID
			parsedTraceID, err := parseHexTraceID(traceID)
			if err == nil {
				// Create new span context
				newSpanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    parsedTraceID,
					SpanID:     trace.SpanID{}, // Will be set by tracer
					TraceFlags: trace.FlagsSampled,
					Remote:     true,
				})
				// Update context with new span context
				ctx = trace.ContextWithSpanContext(ctx, newSpanContext)
				c.Request = c.Request.WithContext(ctx)
			}

			// Set trace ID in gin context and header
			c.Set(GinTraceIDKey, traceID)
			ctx := context.WithValue(c.Request.Context(), GinTraceIDKey, traceID)
			c.Request = c.Request.WithContext(ctx)
			c.Header(TraceIDHeader, traceID)
		}

		c.Next()
	}
}

// GetTraceID retrieves trace ID from gin context
func GetTraceID(c *gin.Context) (string, bool) {
	if v, exists := c.Get(GinTraceIDKey); exists {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}
