package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// ZapGormLogger is a logger that implements the GORM logger interface using zap
type ZapGormLogger struct {
	logger *zap.Logger
}

// NewZapGormLogger creates a new ZapGormLogger
func NewZapGormLogger(zapLogger *zap.Logger, level logger.LogLevel) *ZapGormLogger {
	return &ZapGormLogger{
		logger: zapLogger,
	}
}

// LogMode sets the log level
func (l *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	return &newLogger
}

// Info logs an info message
func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...any) {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Infow(msg, data...)
		return
	}
	l.logger.Sugar().Infow(msg, data...)
}

// Warn logs a warning message
func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Warnw(msg, data...)
		return
	}
	l.logger.Sugar().Warnw(msg, data...)
}

// Error logs an error message
func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...any) {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Errorw(msg, data...)
		return
	}
	l.logger.Sugar().Errorw(msg, data...)
}

// Trace logs a trace message
func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
			l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Errorw("trace",
				"err", err,
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
			return
		}
		l.logger.Sugar().Errorw("trace",
			"err", err,
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	case elapsed > 200*time.Millisecond:
		sql, rows := fc()
		if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
			l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Warnw("trace",
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
			return
		}
		l.logger.Sugar().Warnw("trace",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	default:
		sql, rows := fc()
		if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
			l.logger.With(zap.String(traceIDKey, traceID)).Sugar().Infow("trace",
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			)
			return
		}
		l.logger.Sugar().Infow("trace",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	}
}
