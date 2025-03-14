package log

import (
	"os"
	"testing"

	"github.com/fize/go-ext/config"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.LogConfig
		wantErr bool
	}{
		{
			name: "valid console logger",
			cfg: &config.LogConfig{
				Level:  "info",
				Format: "string",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "valid json file logger",
			cfg: &config.LogConfig{
				Level:      "debug",
				Format:     "json",
				Output:     "file",
				Filename:   "test.log",
				MaxSize:    10,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
		{
			name: "invalid level",
			cfg: &config.LogConfig{
				Level:  "invalid",
				Format: "string",
				Output: "stdout",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := InitLogger(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.NotNil(t, logger.sugar)
			}
		})
	}

	// Clean up test log file
	os.Remove("test.log")
}

func TestLoggerMethods(t *testing.T) {
	logger, err := InitLogger(DefaultConfig())
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Test all logging methods
	tests := []struct {
		name  string
		logFn func()
	}{
		{
			name:  "Debug",
			logFn: func() { logger.Debug("debug message") },
		},
		{
			name:  "Debugf",
			logFn: func() { logger.Debugf("debug %s", "message") },
		},
		{
			name:  "Debugw",
			logFn: func() { logger.Debugw("debug message", "key", "value") },
		},
		// Add similar tests for Info, Warn, Error levels
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simply verify that the logging methods don't panic
			assert.NotPanics(t, tt.logFn)
		})
	}
}
