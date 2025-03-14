// This file is used to configure the server
package config

import (
	"fmt"
	"net"
	"regexp"
)

// Default values for server configuration
const (
	_defaultBindAddr = "0.0.0.0:8080"
)

// Metrics defines the metrics configuration options
type Metrics struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"serviceName"`
	Path        string `mapstructure:"path"`
	Port        int    `mapstructure:"port"`
	// TimeSensitive is a flag to enable time sensitive metrics buckets
	TimeSensitive bool `mapstructure:"timeSensitive"`
	// ExcludeItem is a list of items to exclude from metrics, e.g. /health
	ExcludeItem []string `mapstructure:"excludeItem"`
}

// Trace defines the trace configuration options
type Trace struct {
	Enabled bool `mapstructure:"enabled"`
	// ServiceName is the application service name
	ServiceName string `mapstructure:"serviceName"`
	// Stdout is a flag to enable/disable stdout
	Stdout bool `mapstructure:"stdout"`
	// Endpoint is the OTEL collector endpoint to send traces to
	Endpoint string `mapstructure:"endpoint"`
	// ExcludeItem is a list of items to exclude from tracing, e.g. /health
	ExcludeItem []string `mapstructure:"excludeItem"`
}

// ServerConfig defines the server configuration options
type ServerConfig struct {
	// Server bind address, e.g. 0.0.0.0:8080
	BindAddr string `mapstructure:"bindAddr"`
	// Metrics
	Metrics *Metrics `mapstructure:"metrics"`
	// Trace
	Trace *Trace `mapstructure:"trace"`
}

// ServerConfigOption is used to configure the server
type ServerConfigOption func(*ServerConfig)

// defaultServerConfig returns the default server configuration
func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr: _defaultBindAddr,
		Metrics: &Metrics{
			Enabled:     false,
			ServiceName: "default",
		},
		Trace: &Trace{
			Enabled:     false,
			ServiceName: "default",
		},
	}
}

// NewServerConfig creates a new ServerConfig with the given options
func NewServerConfig(opts ...ServerConfigOption) (*ServerConfig, error) {
	cfg := defaultServerConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	if err := validateAddr(cfg.BindAddr); err != nil {
		return nil, err
	}

	if cfg.Metrics != nil {
		if cfg.Metrics.Enabled {
			if err := validateMetricsServiceName(cfg.Metrics.ServiceName); err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}

// WithBindAddr sets the server bind address
func WithBindAddr(addr string) ServerConfigOption {
	return func(c *ServerConfig) {
		c.BindAddr = addr
	}
}

// WithMetrics sets the metrics configuration
func WithMetrics(metrics *Metrics) ServerConfigOption {
	return func(c *ServerConfig) {
		c.Metrics = metrics
	}
}

// WithMetricsEnabled enables/disables metrics
func WithMetricsEnabled(enabled bool) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.Enabled = enabled
	}
}

// WithMetricsPath sets the metrics path
func WithMetricsPath(path string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.Path = path
	}
}

// WithMetricsPort sets the metrics port
func WithMetricsPort(port int) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.Port = port
	}
}

// WithMetricsServiceName sets the metrics service name
func WithMetricsServiceName(name string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.ServiceName = name
	}
}

func WithMetricsExcludeItem(items []string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.ExcludeItem = items
	}
}

func WithMetricsTimeSensitive(enabled bool) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Metrics == nil {
			c.Metrics = &Metrics{}
		}
		c.Metrics.TimeSensitive = enabled
	}
}

func WithTrace(trace *Trace) ServerConfigOption {
	return func(c *ServerConfig) {
		c.Trace = trace
	}
}

func WithTraceEnabled(enabled bool) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Trace == nil {
			c.Trace = &Trace{}
		}
		c.Trace.Enabled = enabled
	}
}

func WithTraceServiceName(name string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Trace == nil {
			c.Trace = &Trace{}
		}
		c.Trace.ServiceName = name
	}
}

func WithTraceStdout(enabled bool) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Trace == nil {
			c.Trace = &Trace{}
		}
		c.Trace.Stdout = enabled
	}
}

func WithTraceEndpoint(endpoint string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Trace == nil {
			c.Trace = &Trace{}
		}
		c.Trace.Endpoint = endpoint
	}
}

func WithTraceExcludeItem(items []string) ServerConfigOption {
	return func(c *ServerConfig) {
		if c.Trace == nil {
			c.Trace = &Trace{}
		}
		c.Trace.ExcludeItem = items
	}
}

func validateAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address format %s: %v", addr, err)
	}
	// Check if host is valid IP
	if host != "" && host != "0.0.0.0" && host != "localhost" {
		if ip := net.ParseIP(host); ip == nil {
			return fmt.Errorf("invalid IP address: %s", host)
		}
	}
	// Check if port is valid
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("invalid port number: %s", port)
	}
	return nil
}

func validateMetricsServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name is empty")
	} else {
		r := regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
		if !r.MatchString(name) {
			return fmt.Errorf("serviceName '%s' is invalid, must match ^[a-zA-Z_:][a-zA-Z0-9_:]*$",
				name)
		}
	}
	return nil
}
