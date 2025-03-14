package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestTraceIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TraceID())

	r.GET("/test", func(c *gin.Context) {
		traceID, exists := GetTraceID(c)
		assert.True(t, exists)
		assert.Len(t, traceID, 32) // hex encoded 16 bytes
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name        string
		headerName  string
		headerValue string
	}{
		{
			name: "No TraceID Header",
		},
		{
			name:        "With TraceID Header",
			headerName:  TraceIDHeader,
			headerValue: "1234567890abcdef1234567890abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.headerName != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			traceIDHeader := w.Header().Get(TraceIDHeader)
			assert.NotEmpty(t, traceIDHeader)

			if tt.headerValue != "" {
				assert.Equal(t, tt.headerValue, traceIDHeader)
			}
		})
	}
}

func TestGenerateAndParseTraceID(t *testing.T) {
	traceID := generateTraceID()
	assert.Len(t, traceID, 32)

	parsed, err := parseHexTraceID(traceID)
	assert.NoError(t, err)
	assert.NotEqual(t, trace.TraceID{}, parsed)
}

func TestFromContext(t *testing.T) {
	tests := []struct {
		name     string
		traceID  string
		wantID   string
		wantBool bool
	}{
		{
			name:     "With TraceID",
			traceID:  "1234567890abcdef1234567890abcdef",
			wantID:   "1234567890abcdef1234567890abcdef",
			wantBool: true,
		},
		{
			name:     "Without TraceID",
			wantID:   "",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.traceID != "" {
				ctx = context.WithValue(ctx, GinTraceIDKey, tt.traceID)
			}

			gotID, gotBool := FromContext(ctx)
			assert.Equal(t, tt.wantID, gotID)
			assert.Equal(t, tt.wantBool, gotBool)
		})
	}
}
