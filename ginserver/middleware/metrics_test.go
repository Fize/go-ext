package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fize/go-ext/config"
	"github.com/gin-gonic/gin"
)

func TestMetricsMiddleware(t *testing.T) {
	cfg := &config.Metrics{
		Enabled:     true,
		ServiceName: "test_service",
		Path:        "/metrics",
		Port:        9090,
	}

	r := gin.Default()
	r.Use(MetricsMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
