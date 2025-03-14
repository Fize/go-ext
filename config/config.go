// config package is a package for managing application configuration.
// It provides a BaseConfig struct that holds the configuration for the application.
// The package also provides functions to load the configuration from a file or environment variables, set default values, and validate the configuration.
package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Email
type Email struct {
	Account  string `fig:"account"`
	SMTP     string `fig:"smtp"`
	Port     int    `fig:"port"`
	Password string `fig:"password"`
}

var once sync.Once

// NewConfig creates a new configuration with default values
func NewConfig() *BaseConfig {
	return &BaseConfig{
		Server: defaultServerConfig(),
		SQL:    defaultSQLConfig(),
		Log:    defaultLogConfig(),
	}
}

// BaseConfig is the base configuration
type BaseConfig struct {
	// server configuration
	Server *ServerConfig `mapstructure:"server"`
	// SQL database configuration
	SQL *SQLConfig `mapstructure:"sql"`
	// log configuration
	Log *LogConfig `mapstructure:"log"`
	// custom configuration, can be any type
	Custom interface{} `mapstructure:"custom"`
}

// WithCustomConfig sets the custom configuration
func WithCustomConfig(custom interface{}) func(*BaseConfig) {
	return func(bc *BaseConfig) {
		bc.Custom = custom
	}
}

// ParseCustomConfig parses the custom configuration into the provided interface
func (bc *BaseConfig) ParseCustomConfig(out interface{}) error {
	if bc.Custom == nil {
		return nil
	}

	// Use mapstructure to convert the custom config to the desired type
	config := &mapstructure.DecoderConfig{
		Result:           out,
		WeaklyTypedInput: true,
		TagName:          "mapstructure",
		// Disable case insensitive matching
		MatchName: func(mapKey, fieldName string) bool {
			return mapKey == fieldName
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %v", err)
	}

	return decoder.Decode(bc.Custom)
}

// Load loads configuration from the specified file and optionally parses custom config
// name: configuration file path
// env: whether to load from environment variables
func (bc *BaseConfig) Load(name string, env bool) {
	once.Do(func() {
		v := viper.New()

		// Configure Viper to preserve key case sensitivity
		v.SetEnvPrefix("ext")
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AllowEmptyEnv(true)

		if !env {
			v.SetConfigFile(name)
			if err := v.ReadInConfig(); err != nil {
				panic(fmt.Errorf("failed to read config file: %v", err))
			}
		}

		// Set default values
		bc.setDefaults(v)

		// Use Viper's Unmarshal to parse the configuration
		if err := v.Unmarshal(bc, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "mapstructure"
		}); err != nil {
			panic(fmt.Errorf("failed to unmarshal config: %v", err))
		}

		// Initialize and validate configuration
		if err := bc.initAndValidate(); err != nil {
			panic(err)
		}
	})
}

// setDefaults sets default values for the configuration
func (bc *BaseConfig) setDefaults(v *viper.Viper) {
	// Set default database configuration
	v.SetDefault("sql.type", defaultSQLConfig().Type)
	v.SetDefault("sql.host", defaultSQLConfig().Host)
	v.SetDefault("sql.user", defaultSQLConfig().User)
	v.SetDefault("sql.password", defaultSQLConfig().Password)
	v.SetDefault("sql.db", defaultSQLConfig().DB)
	v.SetDefault("sql.maxIdleConns", defaultSQLConfig().MaxIdleConns)
	v.SetDefault("sql.maxOpenConns", defaultSQLConfig().MaxOpenConns)
	v.SetDefault("sql.debug", defaultSQLConfig().Debug)

	// Set default log configuration
	v.SetDefault("log.filename", defaultLogConfig().Filename)
	v.SetDefault("log.maxSize", defaultLogConfig().MaxSize)
	v.SetDefault("log.maxBackups", defaultLogConfig().MaxBackups)
	v.SetDefault("log.maxAge", defaultLogConfig().MaxAge)
	v.SetDefault("log.compress", defaultLogConfig().Compress)
	v.SetDefault("log.level", defaultLogConfig().Level)
	v.SetDefault("log.format", defaultLogConfig().Format)
	v.SetDefault("log.output", defaultLogConfig().Output)

	// Set default server configuration
	v.SetDefault("server.bindAddr", defaultServerConfig().BindAddr)

	// Set default ServerMetrics
	v.SetDefault("server.metrics.path", defaultServerConfig().Metrics.Path)
	v.SetDefault("server.metrics.port", defaultServerConfig().Metrics.Port)
	v.SetDefault("server.metrics.enabled", defaultServerConfig().Metrics.Enabled)
	v.SetDefault("server.metrics.serviceName", defaultServerConfig().Metrics.ServiceName)
	v.SetDefault("server.metrics.excludeItem", defaultServerConfig().Metrics.ExcludeItem)
	v.SetDefault("server.metrics.timeSensitive", defaultServerConfig().Metrics.TimeSensitive)

	// Set default ServerTrace
	v.SetDefault("server.trace.enabled", defaultServerConfig().Trace.Enabled)
	v.SetDefault("server.trace.serviceName", defaultServerConfig().Trace.ServiceName)
	v.SetDefault("server.trace.stdout", defaultServerConfig().Trace.Stdout)
	v.SetDefault("server.trace.endpoint", defaultServerConfig().Trace.Endpoint)
	v.SetDefault("server.trace.excludeItem", defaultServerConfig().Trace.ExcludeItem)
}

// initAndValidate initializes and validates the configuration
func (bc *BaseConfig) initAndValidate() error {
	// Validate SQL config
	sqlCfg, err := NewSQLConfig(
		WithType(bc.SQL.Type),
		WithHost(bc.SQL.Host),
		WithUser(bc.SQL.User),
		WithPassword(bc.SQL.Password),
		WithDB(bc.SQL.DB),
		WithMaxIdleConns(bc.SQL.MaxIdleConns),
		WithMaxOpenConns(bc.SQL.MaxOpenConns),
		WithDebug(bc.SQL.Debug),
	)
	if err != nil {
		return fmt.Errorf("invalid SQL config: %v", err)
	}
	bc.SQL = sqlCfg

	// Validate Log config
	logCfg, err := NewLogConfig(
		WithFilename(bc.Log.Filename),
		WithMaxSize(bc.Log.MaxSize),
		WithMaxBackups(bc.Log.MaxBackups),
		WithMaxAge(bc.Log.MaxAge),
		WithCompress(bc.Log.Compress),
		WithLevel(bc.Log.Level),
		WithFormat(bc.Log.Format),
		WithOutput(bc.Log.Output),
	)
	if err != nil {
		return fmt.Errorf("invalid Log config: %v", err)
	}
	bc.Log = logCfg

	serCfg, err := NewServerConfig(
		WithBindAddr(bc.Server.BindAddr),
		WithMetrics(bc.Server.Metrics),
		WithMetricsEnabled(bc.Server.Metrics.Enabled),
		WithMetricsPath(bc.Server.Metrics.Path),
		WithMetricsPort(bc.Server.Metrics.Port),
		WithMetricsServiceName(bc.Server.Metrics.ServiceName),
		WithMetricsExcludeItem(bc.Server.Metrics.ExcludeItem),
		WithMetricsTimeSensitive(bc.Server.Metrics.TimeSensitive),
		WithTrace(bc.Server.Trace),
		WithTraceEnabled(bc.Server.Trace.Enabled),
		WithTraceServiceName(bc.Server.Trace.ServiceName),
		WithTraceStdout(bc.Server.Trace.Stdout),
		WithTraceEndpoint(bc.Server.Trace.Endpoint),
		WithTraceExcludeItem(bc.Server.Trace.ExcludeItem),
	)
	if err != nil {
		return fmt.Errorf("invalid Server config: %v", err)
	}
	bc.Server = serCfg

	return nil
}
