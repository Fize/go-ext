package middleware

import (
	"context"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/fize/go-ext/config"
	"github.com/fize/go-ext/log"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	tracer    oteltrace.Tracer
	traceOnce = sync.Once{}

	exTracePathPatterns []*regexp.Regexp
	ip                  string
	hostname            = ""
)

func initInfo() {
	if ip == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Errorf("failed to get ip: %v", err)
			return
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Errorf("failed to get hostname: %v", err)
	}
}

// GetTracer returns the tracer for the application
func GetTracer() oteltrace.Tracer {
	return tracer
}

// SetTracer sets the tracer for the application
func SetTracer(name string) {
	traceOnce.Do(func() {
		if tracer != nil {
			return
		}
		tracer = otel.Tracer(name)
	})
}

func InitTracer(ctx context.Context, cfg *config.Trace) (*sdktrace.TracerProvider, error) {
	initTraceExcludePath(cfg.ExcludeItem)
	initInfo()
	var err error
	var exporter sdktrace.SpanExporter
	if cfg.Stdout {
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stdout),
			stdouttrace.WithoutTimestamps(),
		)
	} else {
		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(cfg.Endpoint), otlptracegrpc.WithInsecure())
	}
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("service.name", cfg.ServiceName),
			semconv.ServiceInstanceIDKey.String(ip),
			semconv.HostNameKey.String(hostname),
			attribute.String("k8s.pod.ip", ip),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func initTraceExcludePath(exi []string) {
	if len(exi) > 0 {
		for _, item := range exi {
			// Compile regular expression pattern
			pattern, err := regexp.Compile(item)
			if err != nil {
				panic(err)
			}
			exTracePathPatterns = append(exTracePathPatterns, pattern)
		}
	}
}

// isTracePathExcluded checks if the given path matches any exclude pattern
func isTracePathExcluded(path string) bool {
	for _, pattern := range exTracePathPatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

// TraceFilter returns a function that filters traces based on the configuration
func TraceFilter(cfg *config.Trace) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		return !isTracePathExcluded(c.Request.URL.Path)
	}
}
