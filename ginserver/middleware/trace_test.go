package middleware

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/fize/go-ext/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitTracer(t *testing.T) {
	cfg := &config.Trace{
		Enabled:     true,
		ServiceName: "test-service",
		Stdout:      true,
		ExcludeItem: []string{"^/health$", "^/metrics"},
	}

	tp, err := InitTracer(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, tp)
	assert.NotEmpty(t, exTracePathPatterns)
}

func TestTraceFilter(t *testing.T) {
	// Reset patterns for clean test
	exTracePathPatterns = nil

	cfg := &config.Trace{
		ExcludeItem: []string{
			"^/health$",       // 精确匹配 /health
			"^/metrics$",      // 精确匹配 /metrics
			"^/debug/pprof.*", // 匹配所有 pprof 路径
			".+\\.jpg$",       // 匹配所有 .jpg 结尾的路径
		},
	}
	initTraceExcludePath(cfg.ExcludeItem)

	tests := []struct {
		path     string
		expected bool
	}{
		{"/health", false},           // 被排除
		{"/healthz", true},           // 不完全匹配 /health
		{"/metrics", false},          // 被排除
		{"/metrics/custom", true},    // 不完全匹配 /metrics
		{"/api/user", true},          // 正常路径
		{"/image.jpg", false},        // 被排除
		{"/debug/pprof/", false},     // 被排除
		{"/debug/pprof/heap", false}, // 被排除
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request = &http.Request{
				URL: &url.URL{Path: tt.path},
			}
			filter := TraceFilter(cfg)
			result := filter(c)
			assert.Equal(t, tt.expected, result, "Path: %s should return %v", tt.path, tt.expected)
		})
	}
}

func TestIsTracePathExcluded(t *testing.T) {
	// Reset patterns
	exTracePathPatterns = nil

	patterns := []string{
		"^/health$",
		"^/metrics",
		".+\\.jpg$",
	}
	initTraceExcludePath(patterns)

	tests := []struct {
		path     string
		expected bool
	}{
		{"/health", true},
		{"/healthz", false},
		{"/metrics", true},
		{"/metrics/custom", true},
		{"/api/user", false},
		{"/image.jpg", true},
		{"/path/file.jpg", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isTracePathExcluded(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
