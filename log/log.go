package log

import (
	"os"

	"github.com/fize/go-ext/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	logLevel      = zap.NewAtomicLevel()
)

// logger := log.InitLogger()
// defer logger.Sync()
func InitLogger() *zap.Logger {
	setLevel()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugaredLogger = logger.Sugar()
	return logger
}

func Sync() {
	sugaredLogger.Sync()
}

func setLevel() {
	switch config.Read().Log.Level {
	case "debug":
		logLevel.SetLevel(zapcore.Level(zapcore.DebugLevel))
	case "info":
		logLevel.SetLevel(zapcore.Level(zapcore.InfoLevel))
	case "warn":
		logLevel.SetLevel(zapcore.Level(zapcore.WarnLevel))
	case "error":
		logLevel.SetLevel(zapcore.Level(zapcore.ErrorLevel))
	default:
		logLevel.SetLevel(zapcore.Level(zapcore.DebugLevel))
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if config.Read().Log.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	if config.Read() != nil {
		if config.Read().Log.Output == "file" {
			lumberJackLogger := &lumberjack.Logger{
				Filename:   config.Read().Log.Filename,
				MaxSize:    config.Read().Log.MaxSize,
				MaxBackups: config.Read().Log.MaxBackups,
				MaxAge:     config.Read().Log.MaxAge,
				Compress:   config.Read().Log.Compress,
			}
			return zapcore.AddSync(lumberJackLogger)
		}
	}
	return zapcore.AddSync(zapcore.Lock(os.Stdout))
}

func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

func Debugw(msg string, args ...interface{}) {
	sugaredLogger.Debugw(msg, args...)
}

func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

func Infow(msg string, args ...interface{}) {
	sugaredLogger.Infow(msg, args...)
}

func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

func Warnw(msg string, args ...interface{}) {
	sugaredLogger.Warnw(msg, args...)
}

func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

func Errorw(err error, args ...interface{}) {
	sugaredLogger.Errorw(err.Error(), args...)
}

func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}

func Fatalw(msg string, args ...interface{}) {
	sugaredLogger.Fatalw(msg, args...)
}
