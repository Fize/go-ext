// Description: This file is used to collect metrics for the service.
package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/fize/go-ext/config"
	"github.com/fize/go-ext/log"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
)

// Meter returns the metrics meter used for custom metrics.
func Meter() api.Meter {
	return meter
}

var meter api.Meter

var (
	cfg *config.Metrics

	// Change exPath to store compiled regular expressions
	exPathPatterns []*regexp.Regexp
)

const (
	// Called, generally means that the current service is invoked by another service
	Callee = "callee"
)

var (
	// if data type is counter, otlp will add a _total suffix to the metric name,
	// so we use updowncounter to avoid this suffix
	RequestCount    api.Float64UpDownCounter
	ErrRequestCount api.Float64UpDownCounter
	RequestTime     api.Float64Histogram
)

var (
	normalBuckets        = []float64{5, 10, 25, 50, 100, 150, 200, 250, 500, 2000}
	timeSensitiveBuckets = []float64{1, 3, 5, 7, 10, 15, 40, 80, 200, 500}
)

// initMetrics initializes the metrics
func initMetrics(c *config.Metrics) {
	cfg = c
	var err error
	var sr, se, srt string
	sr = fmt.Sprintf("%s_server_requests_seconds_total", cfg.ServiceName)
	se = fmt.Sprintf("%s_server_requests_error_total", cfg.ServiceName)
	srt = fmt.Sprintf("%s_server_requests_time", cfg.ServiceName)
	meter = initOtelMetrics(srt)
	RequestCount, err = meter.Float64UpDownCounter(sr, api.WithDescription("count of server requests"))
	if err != nil {
		log.Fatalf("failed to initialize server_requests_seconds_total %v", err)
	}
	ErrRequestCount, err = meter.Float64UpDownCounter(se, api.WithDescription("count of server request error"))
	if err != nil {
		log.Fatalf("failed to initialize server_requests_error_total %v", err)
	}
	RequestTime, err = meter.Float64Histogram(srt, api.WithDescription("server request latency distributions. ms"))
	if err != nil {
		log.Fatalf("failed to initialize server_requests_time %v", err)
	}
	initExcludePath()
	go metricsServer()
}

func metricsServer() {
	log.Debugf("serving metrics at :%d%s", cfg.Port, cfg.Path)
	http.Handle(cfg.Path, promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		log.Fatalf("failed to serve metrics: %v", err)
	}
}

func initOtelMetrics(m string) api.Meter {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("failed to initialize prometheus exporter %v", err)
	}
	var ms metric.Stream
	if cfg.TimeSensitive {
		ms = metric.Stream{
			Name: m,
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: timeSensitiveBuckets,
			},
		}
	} else {
		ms = metric.Stream{
			Name: m,
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: normalBuckets,
			},
		}
	}
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithView(metric.NewView(
			metric.Instrument{
				Name:  m,
				Scope: instrumentation.Scope{Name: cfg.ServiceName},
			},
			ms,
		)),
	)
	return provider.Meter(cfg.ServiceName)
}

func initExcludePath() {
	if len(cfg.ExcludeItem) > 0 {
		for _, item := range cfg.ExcludeItem {
			// Compile regular expression pattern
			pattern, err := regexp.Compile(item)
			if err != nil {
				panic(err)
			}
			exPathPatterns = append(exPathPatterns, pattern)
		}
	}
}

// isPathExcluded checks if the given path matches any exclude pattern
func isPathExcluded(path string) bool {
	for _, pattern := range exPathPatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

// MetricsMiddleware is a middleware that collects metrics
func MetricsMiddleware(cfg *config.Metrics) func(c *gin.Context) {
	initMetrics(cfg)
	return func(c *gin.Context) {
		// Use the new isPathExcluded function
		if isPathExcluded(c.Request.URL.EscapedPath()) {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		ginMetricHandle(c, start)
	}
}

func ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer
	path := ctx.FullPath()
	method := r.Method
	statusCode := w.Status()
	opt := setLabel(cfg.ServiceName, path, method, statusCode)
	RequestCount.Add(ctx, 1, opt)
	latency := time.Since(start)
	RequestTime.Record(ctx, float64(latency.Milliseconds()), opt)
	if statusCode >= 400 {
		ErrRequestCount.Add(ctx, 1, opt)
	}
}

func setLabel(name, path, method string, code int) api.MeasurementOption {
	return api.WithAttributes(
		attribute.Key("application").String(name),
		attribute.Key("interface").String(path),
		attribute.Key("call_type").String(Callee),
		attribute.Key("method").String(method),
		attribute.Key("error_code").Int(code),
	)
}
