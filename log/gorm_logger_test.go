package log

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

func setupLogger() (*zap.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := zap.New(core)
	return logger, &buf
}

func TestZapGormLogger_Info(t *testing.T) {
	zapLogger, buf := setupLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Info)

	ctx := context.WithValue(context.Background(), traceIDKey, "123456")
	gormLogger.Info(ctx, "info message", "key", "value")

	assert.Contains(t, buf.String(), `"msg":"info message"`)
	assert.Contains(t, buf.String(), `"trace-id":"123456"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}

func TestZapGormLogger_Warn(t *testing.T) {
	zapLogger, buf := setupLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Warn)

	ctx := context.WithValue(context.Background(), traceIDKey, "123456")
	gormLogger.Warn(ctx, "warn message", "key", "value")

	assert.Contains(t, buf.String(), `"msg":"warn message"`)
	assert.Contains(t, buf.String(), `"trace-id":"123456"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}

func TestZapGormLogger_Error(t *testing.T) {
	zapLogger, buf := setupLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Error)

	ctx := context.WithValue(context.Background(), traceIDKey, "123456")
	gormLogger.Error(ctx, "error message", "key", "value")

	assert.Contains(t, buf.String(), `"msg":"error message"`)
	assert.Contains(t, buf.String(), `"trace-id":"123456"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}

func TestZapGormLogger_Trace(t *testing.T) {
	zapLogger, buf := setupLogger()
	gormLogger := NewZapGormLogger(zapLogger, logger.Info)

	ctx := context.WithValue(context.Background(), traceIDKey, "123456")
	begin := time.Now()
	fc := func() (string, int64) {
		return "SELECT * FROM users", 10
	}
	gormLogger.Trace(ctx, begin, fc, nil)

	assert.Contains(t, buf.String(), `"msg":"trace"`)
	assert.Contains(t, buf.String(), `"trace-id":"123456"`)
	assert.Contains(t, buf.String(), `"sql":"SELECT * FROM users"`)
	assert.Contains(t, buf.String(), `"rows":10`)
}
