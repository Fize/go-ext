// This file is used to configure the SQL-database
package config

import (
	"fmt"
)

// default configuration
const (
	// default database type
	_defaultSQLType = "sqlite3"
	// default database file
	_defaultSQL = "./sqlite.db"
)

// support mysql and sqlite
const (
	MySQL   = "mysql"
	Sqlite3 = "sqlite3"
)

// SQLConfig is used to configure the SQL-database
type SQLConfig struct {
	// Database type only support mysql and sqlite, default sqlite
	Type string `mapstructure:"type"`
	// Database host, include port such as 127.0.0.1:3306
	Host string `mapstructure:"host"`
	// Database user
	User string `mapstructure:"user"`
	// Database password
	Password string `mapstructure:"password"`
	// Database name
	DB string `mapstructure:"db"`
	// Maximum number of idle connections
	MaxIdleConns int `mapstructure:"maxIdleConns"`
	// Maximum number of open connections
	MaxOpenConns int `mapstructure:"maxOpenConns"`
	// Print raw sql for debugging
	Debug bool `mapstructure:"debug"`
}

// SQLConfigOption is used to configure the SQL-database
type SQLConfigOption func(*SQLConfig)

func defaultSQLConfig() *SQLConfig {
	return &SQLConfig{
		Type: _defaultSQLType,
		DB:   _defaultSQL,
	}
}

// NewSQLConfig creates a new SQLConfig with the given options
func NewSQLConfig(opts ...SQLConfigOption) (*SQLConfig, error) {
	cfg := defaultSQLConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	// Validate database type
	if cfg.Type != MySQL && cfg.Type != Sqlite3 {
		return nil, fmt.Errorf("invalid database type: %s", cfg.Type)
	}

	return cfg, nil
}

// WithType sets the database type
func WithType(dbType string) SQLConfigOption {
	return func(c *SQLConfig) {
		c.Type = dbType
	}
}

// WithHost sets the database host
func WithHost(host string) SQLConfigOption {
	return func(c *SQLConfig) {
		c.Host = host
	}
}

// WithUser sets the database user
func WithUser(user string) SQLConfigOption {
	return func(c *SQLConfig) {
		c.User = user
	}
}

// WithPassword sets the database password
func WithPassword(password string) SQLConfigOption {
	return func(c *SQLConfig) {
		c.Password = password
	}
}

// WithDB sets the database name
func WithDB(db string) SQLConfigOption {
	return func(c *SQLConfig) {
		c.DB = db
	}
}

// WithMaxIdleConns sets the maximum number of idle connections
func WithMaxIdleConns(n int) SQLConfigOption {
	return func(c *SQLConfig) {
		c.MaxIdleConns = n
	}
}

// WithMaxOpenConns sets the maximum number of open connections
func WithMaxOpenConns(n int) SQLConfigOption {
	return func(c *SQLConfig) {
		c.MaxOpenConns = n
	}
}

// WithDebug sets the SQL debug mode
func WithDebug(debug bool) SQLConfigOption {
	return func(c *SQLConfig) {
		c.Debug = debug
	}
}
