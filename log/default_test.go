package log

import (
	"fmt"
	"testing"

	"github.com/fize/go-ext/config"
	"github.com/stretchr/testify/assert"
)

func testinit() {
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

func TestDefaultLogger(t *testing.T) {
	testinit()
	defer Sync()
	// Test that the default logger is initialized
	assert.NotNil(t, defaultLogger)

	// Test global logging functions
	tests := []struct {
		name  string
		logFn func()
	}{
		{
			name:  "Debug",
			logFn: func() { Debugw("debug message") },
		},
		{
			name:  "Info",
			logFn: func() { Info("info message") },
		},
		{
			name:  "Warn",
			logFn: func() { Warn("warn message") },
		},
		{
			name:  "Error",
			logFn: func() { Error("error message") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, tt.logFn)
		})
	}
}
