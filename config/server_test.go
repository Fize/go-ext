package config

import (
	"os"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	// TestNewServerConfig verifies that creating a new ServerConfig with various options works as expected
	tests := []struct {
		name    string
		opts    []ServerConfigOption
		want    *ServerConfig
		wantErr bool
	}{
		{
			name: "default config",
			opts: nil,
			want: &ServerConfig{
				BindAddr: _defaultBindAddr,
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "valid custom config with metrics",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetrics(&Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/custom-metrics",
					Port:        9090,
				}),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/custom-metrics",
					Port:        9090,
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "valid custom config with metrics and trace",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetrics(&Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/custom-metrics",
					Port:        9090,
				}),
				WithTrace(&Trace{
					Enabled:  true,
					Endpoint: "http://localhost:4317",
				}),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/custom-metrics",
					Port:        9090,
				},
				Trace: &Trace{
					Enabled:  true,
					Endpoint: "http://localhost:4317",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid bind address",
			opts: []ServerConfigOption{
				WithBindAddr("invalid:addr:8080"),
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:99999"),
			},
			wantErr: true,
		},
		{
			name: "empty bind address",
			opts: []ServerConfigOption{
				WithBindAddr(""),
			},
			wantErr: true,
		},
		{
			name: "localhost address",
			opts: []ServerConfigOption{
				WithBindAddr("localhost:8080"),
			},
			want: &ServerConfig{
				BindAddr: "localhost:8080",
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerConfig(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				compareServerConfig(t, got, tt.want)
			}
		})
	}
}

// Helper function to compare ServerConfig objects
func compareServerConfig(t *testing.T, got, want *ServerConfig) {
	t.Helper()

	if got.BindAddr != want.BindAddr {
		t.Errorf("BindAddr = %v, want %v", got.BindAddr, want.BindAddr)
	}

	// Compare Metrics configuration
	if (got.Metrics == nil) != (want.Metrics == nil) {
		t.Errorf("Metrics presence mismatch: got %v, want %v", got.Metrics != nil, want.Metrics != nil)
		return
	}

	if got.Metrics != nil && want.Metrics != nil {
		if got.Metrics.Enabled != want.Metrics.Enabled {
			t.Errorf("Metrics.Enabled = %v, want %v", got.Metrics.Enabled, want.Metrics.Enabled)
		}
		if got.Metrics.ServiceName != want.Metrics.ServiceName {
			t.Errorf("Metrics.ServiceName = %v, want %v", got.Metrics.ServiceName, want.Metrics.ServiceName)
		}
		if got.Metrics.Path != want.Metrics.Path {
			t.Errorf("Metrics.Path = %v, want %v", got.Metrics.Path, want.Metrics.Path)
		}
		if got.Metrics.Port != want.Metrics.Port {
			t.Errorf("Metrics.Port = %v, want %v", got.Metrics.Port, want.Metrics.Port)
		}
	}

	// Compare Trace configuration
	if (got.Trace == nil) != (want.Trace == nil) {
		t.Errorf("Trace presence mismatch: got %v, want %v", got.Trace != nil, want.Trace != nil)
		return
	}

	if got.Trace != nil && want.Trace != nil {
		if got.Trace.Enabled != want.Trace.Enabled {
			t.Errorf("Trace.Enabled = %v, want %v", got.Trace.Enabled, want.Trace.Enabled)
		}
		if got.Trace.Endpoint != want.Trace.Endpoint {
			t.Errorf("Trace.Endpoint = %v, want %v", got.Trace.Endpoint, want.Trace.Endpoint)
		}
	}
}

func TestValidateAddr(t *testing.T) {
	// TestValidateAddr checks correctness of address validation
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{"valid address", "127.0.0.1:8080", false},
		{"valid default address", "0.0.0.0:8080", false},
		{"valid localhost", "localhost:8080", false},
		{"invalid IP", "256.256.256.256:8080", true},
		{"invalid port", "127.0.0.1:99999", true},
		{"missing port", "127.0.0.1:", true},
		{"missing host", ":8080", false},
		{"invalid format", "invalid_address", true},
		{"empty address", "", true},
		{"multiple colons", "127.0.0.1:8080:8081", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAddr(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerConfigWithMetrics(t *testing.T) {
	// TestServerConfigWithMetrics tests ServerConfig integration with metrics configuration
	tests := []struct {
		name    string
		opts    []ServerConfigOption
		want    *ServerConfig
		wantErr bool
	}{
		{
			name: "metrics disabled",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetricsEnabled(false),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "metrics enabled without service name",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetricsEnabled(true),
				WithMetricsServiceName(""),
			},
			wantErr: true, // should fail because service name is required when metrics is enabled
		},
		{
			name: "valid metrics config",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetrics(&Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/metrics",
					Port:        9090,
				}),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     true,
					ServiceName: "test_service",
					Path:        "/metrics",
					Port:        9090,
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid metrics service name",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetrics(&Metrics{
					Enabled:     true,
					ServiceName: "invalid@service",
				}),
			},
			wantErr: true,
		},
		{
			name: "empty metrics service name",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithMetrics(&Metrics{
					Enabled:     true,
					ServiceName: "",
				}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerConfig(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				compareServerConfig(t, got, tt.want)
			}
		})
	}
}

func TestValidateMetricsServiceName(t *testing.T) {
	// TestValidateMetricsServiceName ensures the metrics service name follows the specified format
	tests := []struct {
		name    string
		svcName string
		wantErr bool
	}{
		{"valid service name", "test_service", false},
		{"valid with colon", "app:service", false},
		{"valid with underscore", "app_service", false},
		{"empty name", "", true},
		{"invalid start char", "1service", true},
		{"invalid chars", "test@service", true},
		{"space in name", "test service", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMetricsServiceName(tt.svcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMetricsServiceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Rename TestLoadFromYAMLFile to TestServerConfigFromYAML
func TestServerConfigFromYAML(t *testing.T) {
	// TestServerConfigFromYAML verifies loading and parsing ServerConfig from a YAML file
	cfg := setupTest()
	cfg.Load("testdata/config.yaml", false)

	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"Server.Metrics.Path", cfg.Server.Metrics.Path, "/metrics"},
		{"Server.Metrics.Port", cfg.Server.Metrics.Port, 9090},
		{"Server.Metrics.Enabled", cfg.Server.Metrics.Enabled, true},
		{"Server.Metrics.ServiceName", cfg.Server.Metrics.ServiceName, "test_service"},
		{"Server.Metrics.TimeSensitive", cfg.Server.Metrics.TimeSensitive, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// Rename TestLoadFromEnv to TestServerConfigFromEnv
func TestServerConfigFromEnv(t *testing.T) {
	// TestServerConfigFromEnv checks loading ServerConfig properties from environment variables
	env := map[string]string{
		"EXT_SERVER_METRICS_PATH":          "/custom-metrics",
		"EXT_SERVER_METRICS_PORT":          "9091",
		"EXT_SERVER_METRICS_ENABLED":       "true",
		"EXT_SERVER_METRICS_SERVICENAME":   "test_service",
		"EXT_SERVER_METRICS_TIMESENSITIVE": "true",
	}

	for k, v := range env {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range env {
			os.Unsetenv(k)
		}
	}()

	cfg := setupTest()
	cfg.Load("", true)

	metricsTests := []struct {
		name     string
		got      any
		expected any
	}{
		{"Metrics.Path", cfg.Server.Metrics.Path, "/custom-metrics"},
		{"Metrics.Port", cfg.Server.Metrics.Port, 9091},
		{"Metrics.Enabled", cfg.Server.Metrics.Enabled, true},
		{"Metrics.ServiceName", cfg.Server.Metrics.ServiceName, "test_service"},
		{"Metrics.TimeSensitive", cfg.Server.Metrics.TimeSensitive, true},
	}

	for _, tt := range metricsTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected Server.%s to be '%v', got '%v'", tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestServerConfigWithTrace(t *testing.T) {
	// TestServerConfigWithTrace tests ServerConfig integration with trace configuration
	tests := []struct {
		name    string
		opts    []ServerConfigOption
		want    *ServerConfig
		wantErr bool
	}{
		{
			name: "trace disabled",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithTraceEnabled(false),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "trace enabled with endpoint",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithTraceEnabled(true),
				WithTraceEndpoint("http://localhost:4317"),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled:  true,
					Endpoint: "http://localhost:4317",
				},
			},
			wantErr: false,
		},
		{
			name: "trace enabled without endpoint",
			opts: []ServerConfigOption{
				WithBindAddr("127.0.0.1:8080"),
				WithTraceEnabled(true),
			},
			want: &ServerConfig{
				BindAddr: "127.0.0.1:8080",
				Metrics: &Metrics{
					Enabled:     false,
					ServiceName: "default",
				},
				Trace: &Trace{
					Enabled: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerConfig(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				compareServerConfig(t, got, tt.want)
			}
		})
	}
}
