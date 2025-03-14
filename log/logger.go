package log

import (
	"fmt"
	"os"

	"github.com/fize/go-ext/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const traceIDKey = "trace-id"

type Logger struct {
	logger     *zap.Logger
	baselogger *zap.Logger
	sugar      *zap.SugaredLogger // rename for clarity
	cfg        *config.LogConfig
}

// InitLogger initializes and returns a new Logger instance
func InitLogger(cfg *config.LogConfig) (*Logger, error) {
	if cfg == nil {
		return nil, fmt.Errorf("log configuration is nil")
	}

	logger := &Logger{
		cfg: cfg,
	}

	if err := logger.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	defaultLogger = logger
	return logger, nil
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() *config.LogConfig {
	return &config.LogConfig{
		Level:  "info",
		Format: "string",
		Output: "stdout",
	}
}

func (l *Logger) GetLogger() *zap.Logger {
	return l.baselogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
}

func (l *Logger) init() error {
	var level zapcore.Level
	switch l.cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		return fmt.Errorf("invalid log level: %s", l.cfg.Level)
	}

	writeSyncer := l.getLogWriter()
	encoder := l.getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level)
	l.baselogger = zap.New(core)
	l.logger = l.baselogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
	l.sugar = l.logger.Sugar() // store the base sugar logger
	return nil
}

func (l *Logger) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if l.cfg.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (l *Logger) getLogWriter() zapcore.WriteSyncer {
	if l.cfg.Output == "file" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   l.cfg.Filename,
			MaxSize:    l.cfg.MaxSize,
			MaxBackups: l.cfg.MaxBackups,
			MaxAge:     l.cfg.MaxAge,
			Compress:   l.cfg.Compress,
		}
		return zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.AddSync(zapcore.Lock(os.Stdout))
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	if err := l.sugar.Sync(); err != nil {
		return err
	}
	return l.logger.Sync()
}

// Logger methods
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Debugw(msg string, args ...interface{}) {
	l.sugar.Debugw(msg, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Infow(msg string, args ...interface{}) {
	l.sugar.Infow(msg, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Warnw(msg string, args ...interface{}) {
	l.sugar.Warnw(msg, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) Errorw(err error, args ...interface{}) {
	l.sugar.Errorw(err.Error(), args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *Logger) Fatalw(msg string, args ...interface{}) {
	l.sugar.Fatalw(msg, args...)
}
