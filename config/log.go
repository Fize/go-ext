// This file is used to configure the log system
package config

import "fmt"

// default configuration
const (
	_defaultLogOutput = "stdout"
	// default log path
	_defaultLogPath = "./zap.log"
	// default log level
	_defaultLogLevel = "info"
	// default log file max size
	_defaultLogMaxSize = 10
	// default log max backups
	_defaultLogMaxBackups = 5
	// default log max age
	_defaultLogMaxAge = 30
	// default log format
	_defaultLogFormat = "string"
)

// Only support string and json
const (
	stringFormat = "string"
	jsonFormat   = "json"
)

// Log levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Log level numbers
const (
	DebugLevelNum = 4
	InfoLevelNum  = 3
	WarnLevelNum  = 2
	ErrorLevelNum = 1
	FatalLevelNum = 0
)

// converts log level string to number
func getLevelNum(level string) (int, error) {
	switch level {
	case DebugLevel:
		return DebugLevelNum, nil
	case InfoLevel:
		return InfoLevelNum, nil
	case WarnLevel:
		return WarnLevelNum, nil
	case ErrorLevel:
		return ErrorLevelNum, nil
	case FatalLevel:
		return FatalLevelNum, nil
	default:
		return -1, fmt.Errorf("unsupported log level: %s", level)
	}
}

// converts log level number to string
func getLevelString(level int) (string, error) {
	// Add range check
	if level < FatalLevelNum || level > 100 {
		return "", fmt.Errorf("unsupported log level number: %d, valid range is [%d-%d]",
			level, FatalLevelNum, DebugLevelNum)
	}

	switch {
	case level >= DebugLevelNum:
		return DebugLevel, nil
	case level == InfoLevelNum:
		return InfoLevel, nil
	case level == WarnLevelNum:
		return WarnLevel, nil
	case level == ErrorLevelNum:
		return ErrorLevel, nil
	case level == FatalLevelNum:
		return FatalLevel, nil
	default:
		if level > InfoLevelNum {
			return DebugLevel, nil
		} else if level > WarnLevelNum {
			return InfoLevel, nil
		} else if level > ErrorLevelNum {
			return WarnLevel, nil
		} else if level > FatalLevelNum {
			return ErrorLevel, nil
		}
		return FatalLevel, nil
	}
}

// LogConfig
type LogConfig struct {
	// log file path
	Filename string `mapstructure:"filename"`
	// log file max size, unit MB
	MaxSize int `mapstructure:"maxSize"`
	// log file max backups
	MaxBackups int `mapstructure:"maxBackups"`
	// log file max age, unit day
	MaxAge int `mapstructure:"maxAge"`
	// log file compress
	Compress bool `mapstructure:"compress"`
	// log level, debug, info, warn, error, fatal
	Level string `mapstructure:"level"`
	// log format, only support string and json
	Format string `mapstructure:"format"`
	// log output, only support stdout and file
	Output string `mapstructure:"output"`
}

// LogConfigConfigOption is used to configure the log system
type LogConfigConfigOption func(*LogConfig)

func defaultLogConfig() *LogConfig {
	return &LogConfig{
		Output:     _defaultLogOutput,
		Filename:   _defaultLogPath,
		MaxSize:    _defaultLogMaxSize,
		MaxBackups: _defaultLogMaxBackups,
		MaxAge:     _defaultLogMaxAge,
		Format:     _defaultLogFormat,
		Level:      _defaultLogLevel,
	}
}

// NewLogConfig creates a new LogConfig config with the given options
func NewLogConfig(opts ...LogConfigConfigOption) (*LogConfig, error) {
	cfg := defaultLogConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	// Validate log format
	if cfg.Format != stringFormat && cfg.Format != jsonFormat {
		return nil, fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	// Validate log level
	if _, err := getLevelNum(cfg.Level); err != nil {
		// Try to parse as number
		levelNum := -1
		if _, err := fmt.Sscanf(cfg.Level, "%d", &levelNum); err == nil {
			levelStr, err := getLevelString(levelNum)
			if err != nil {
				return nil, err
			}
			cfg.Level = levelStr
		} else {
			return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
		}
	}

	return cfg, nil
}

// WithFilename sets the log filename
func WithFilename(filename string) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.Filename = filename
	}
}

// WithMaxSize sets the log file max size
func WithMaxSize(maxSize int) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.MaxSize = maxSize
	}
}

// WithMaxBackups sets the log file max backups
func WithMaxBackups(maxBackups int) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.MaxBackups = maxBackups
	}
}

// WithMaxAge sets the log file max age
func WithMaxAge(maxAge int) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.MaxAge = maxAge
	}
}

// WithCompress sets the log file compress option
func WithCompress(compress bool) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.Compress = compress
	}
}

// WithLevel sets the log level
func WithLevel(level string) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.Level = level
	}
}

// WithFormat sets the log format
func WithFormat(format string) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.Format = format
	}
}

// WithOutput sets the log output
func WithOutput(output string) LogConfigConfigOption {
	return func(c *LogConfig) {
		c.Output = output
	}
}
