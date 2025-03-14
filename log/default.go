package log

import (
	"context"
	"fmt"

	"github.com/fize/go-ext/config"
	"go.uber.org/zap"
)

var defaultLogger *Logger

func init() {
	// Initialize default logger with default configuration
	logger, err := InitLogger(DefaultConfig())
	if err != nil {
		// 如果初始化失败，打印错误并使用 zap 的默认 logger
		fmt.Printf("Failed to initialize default logger: %v\n", err)
		defaultLogger, _ = InitLogger(&config.LogConfig{
			Level:  "info",
			Format: "console",
			Output: "stdout",
		})
	} else {
		defaultLogger = logger
	}
}

// Global functions that use defaultLogger
func Debug(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(args...)
	}
}

func Debugf(template string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugf(template, args...)
	}
}

func Debugw(msg string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugw(msg, args...)
	}
}

func Info(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(args...)
	}
}

func Infof(template string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infof(template, args...)
	}
}

func Infow(msg string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infow(msg, args...)
	}
}

func Warn(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnf(template, args...)
	}
}

func Warnw(msg string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnw(msg, args...)
	}
}

func Error(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorf(template, args...)
	}
}

func Errorw(err error, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorw(err, args...)
	}
}

func Fatal(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(template, args...)
	}
}

func Fatalw(msg string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatalw(msg, args...)
	}
}

// Add other global context-aware functions following the same pattern...

// Sync flushes any buffered log entries from the default logger
func Sync() error {
	if defaultLogger == nil {
		return fmt.Errorf("logger not initialized")
	}
	return defaultLogger.Sync()
}

// clone creates a copy of the Logger with its own sugar logger
func clone() *Logger {
	nl := defaultLogger.baselogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{
		logger:     nl,
		baselogger: defaultLogger.baselogger,
		sugar:      nl.Sugar(),
		cfg:        defaultLogger.cfg,
	}
}

// WithContext returns a new Logger instance with the given context
func WithContext(ctx context.Context) *Logger {
	newLogger := clone()
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		newLogger.sugar = newLogger.sugar.With("traceID", traceID)
	}
	return newLogger
}
