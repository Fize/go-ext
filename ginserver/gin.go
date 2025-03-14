package ginserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fize/go-ext/config"
	"github.com/fize/go-ext/ginserver/middleware"
	"github.com/fize/go-ext/log"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	_ "net/http/pprof"
)

var tp *sdktrace.TracerProvider

const defaultTimeout = 1500 * time.Millisecond

const (
	CORSAllowHeaders = "*"
)

// InitGinServer initializes a new gin server with the given configuration
func InitGinServer(cfg *config.BaseConfig) *gin.Engine {
	r := gin.New()
	r.Use(middleware.TraceID())
	initMetrics(r, cfg.Server.Metrics)
	initTracer(r, cfg.Server.Trace)
	initLoggerAndRecovery(r, cfg.Log)
	return r
}

func initLoggerAndRecovery(r *gin.Engine, cfg *config.LogConfig) {
	logger, err := log.InitLogger(cfg)
	if err != nil {
		panic(err)
	}
	ginlogger := logger.GetLogger()
	if cfg.Level != "debug" && cfg.Level != "info" {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(ginzap.Ginzap(ginlogger, time.RFC3339, true), ginzap.RecoveryWithZap(ginlogger, true))
}

func initMetrics(r *gin.Engine, cfg *config.Metrics) {
	if cfg.Enabled {
		log.Info("setting up metrics middleware")
		r.Use(middleware.MetricsMiddleware(cfg))
	}
}

func initTracer(r *gin.Engine, cfg *config.Trace) {
	if cfg.Enabled {
		var err error
		log.Info("setting up tracing middleware")
		middleware.SetTracer(cfg.ServiceName)
		tp, err = middleware.InitTracer(context.Background(), cfg)
		if err != nil {
			log.Fatalf("init tracer with error: %v", err)
		}
		if cfg.Stdout {
			log.Info("Trace exporter configured for stdout - traces will be printed to terminal")
		}
		r.Use(otelgin.Middleware(cfg.ServiceName, otelgin.WithGinFilter(middleware.TraceFilter(cfg))))
		// r.Use(otelgin.Middleware(cfg.ServiceName))
		log.Info("Tracing middleware successfully initialized")
	} else {
		log.Info("Tracing is disabled in configuration")
	}
}

// HookExit waits for the parent process to exit
// It is used to wait for the parent process to exit
func HookExit(ctx context.Context) {
	defer func() {
		log.Debug("Shutting down tracer provider")
		if tp != nil {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Infof("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	<-ctx.Done()
	log.Info("Parent process exit")
}

// Run starts the gin server with the given configuration
// If gotCtx is true, it will return the context with a cancel function,
// you need to use HookExit to wait for the parent process to exit, HookExit(Run(r, cfg, true)),
// or you can use Run(r, cfg, false) to block the main process
func Run(r *gin.Engine, cfg *config.ServerConfig, gotCtx bool) context.Context {
	srv := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: r.Handler(),
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// pprof server
	go func() {
		log.Info("pprof server started on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof listen: %s\n", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	if gotCtx {
		go gracefulExit(ctx, srv, cancel)
		return ctx
	}
	gracefulExit(ctx, srv, cancel)
	return nil
}

func gracefulExit(ctx context.Context, srv *http.Server, cancel context.CancelFunc) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")
	ctx, cancelChild := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	defer cancelChild()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Infof("timeout of %v.", defaultTimeout)
	}
	log.Info("Server exiting")
}
