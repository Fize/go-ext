package config

import (
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

// setupTest resets the once variable and creates a new configuration
func setupTest() *BaseConfig {
	once = sync.Once{}
	return NewConfig()
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("EXT_SQL_TYPE", "mysql")
	os.Setenv("EXT_SQL_HOST", "localhost:3306")
	os.Setenv("EXT_SQL_USER", "root")
	os.Setenv("EXT_SQL_PASSWORD", "password")
	os.Setenv("EXT_SQL_DB", "testdb")
	os.Setenv("EXT_LOG_LEVEL", "debug")
	os.Setenv("EXT_LOG_FORMAT", "json")
	os.Setenv("EXT_SERVER_TRACE_ENABLED", "true")
	os.Setenv("EXT_SERVER_TRACE_ENDPOINT", "http://localhost:4317")
	defer func() {
		os.Unsetenv("EXT_SQL_TYPE")
		os.Unsetenv("EXT_SQL_HOST")
		os.Unsetenv("EXT_SQL_USER")
		os.Unsetenv("EXT_SQL_PASSWORD")
		os.Unsetenv("EXT_SQL_SQL")
		os.Unsetenv("EXT_LOG_LEVEL")
		os.Unsetenv("EXT_LOG_FORMAT")
		os.Unsetenv("EXT_SERVER_TRACE_ENABLED")
		os.Unsetenv("EXT_SERVER_TRACE_ENDPOINT")
	}()

	once = sync.Once{} // Reset the once variable
	cfg := NewConfig()
	cfg.Load("", true)

	// Validate database configuration
	if cfg.SQL.Type != "mysql" {
		t.Errorf("Expected SQL.Type to be 'mysql', got '%s'", cfg.SQL.Type)
	}
	if cfg.SQL.Host != "localhost:3306" {
		t.Errorf("Expected SQL.Host to be 'localhost:3306', got '%s'", cfg.SQL.Host)
	}
	if cfg.SQL.User != "root" {
		t.Errorf("Expected SQL.User to be 'root', got '%s'", cfg.SQL.User)
	}
	if cfg.SQL.Password != "password" {
		t.Errorf("Expected SQL.Password to be 'password', got '%s'", cfg.SQL.Password)
	}
	if cfg.SQL.DB != "testdb" {
		t.Errorf("Expected SQL.DB to be 'testdb', got '%s'", cfg.SQL.DB)
	}

	// Validate log configuration
	if cfg.Log.Level != "debug" {
		t.Errorf("Expected Log.Level to be 'debug', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Expected Log.Format to be 'json', got '%s'", cfg.Log.Format)
	}

	// Validate trace configuration
	if cfg.Server.Trace.Enabled != true {
		t.Errorf("Expected Server.Trace.Enabled to be 'true', got '%v'", cfg.Server.Trace.Enabled)
	}
	if cfg.Server.Trace.Endpoint != "http://localhost:4317" {
		t.Errorf("Expected Server.Trace.Endpoint to be 'http://localhost:4317', got '%s'", cfg.Server.Trace.Endpoint)
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	cfg := setupTest()
	// Since we're using a .yaml extension, we don't need to explicitly set the config type
	cfg.Load("testdata/config.yaml", false)

	// Validate database configuration
	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"SQL.Type", cfg.SQL.Type, "sqlite3"},
		{"SQL.SQL", cfg.SQL.DB, "./test.db"},
		{"SQL.Debug", cfg.SQL.Debug, true},
		{"Log.Filename", cfg.Log.Filename, "./test.log"},
		{"Log.Level", cfg.Log.Level, "info"},
		{"Log.Format", cfg.Log.Format, "string"},
		{"Server.BindAddr", cfg.Server.BindAddr, "localhost:8080"},
		{"Server.Trace.Enabled", cfg.Server.Trace.Enabled, true},
		{"Server.Trace.Endpoint", cfg.Server.Trace.Endpoint, "http://localhost:4317"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLoadFromJSONFile(t *testing.T) {
	cfg := setupTest()
	viper.SetConfigType("json")
	cfg.Load("testdata/config.json", false)

	// Validate database configuration
	if cfg.SQL.Type != "mysql" {
		t.Errorf("Expected SQL.Type to be 'mysql', got '%s'", cfg.SQL.Type)
	}
	if cfg.SQL.Host != "127.0.0.1:3306" {
		t.Errorf("Expected SQL.Host to be '127.0.0.1:3306', got '%s'", cfg.SQL.Host)
	}
	// ...additional assertions...
}

func TestLoadFromTOMLFile(t *testing.T) {
	cfg := setupTest()
	viper.SetConfigType("toml")
	cfg.Load("testdata/config.toml", false)

	// Validate database configuration
	if cfg.SQL.Type != "sqlite3" {
		t.Errorf("Expected SQL.Type to be 'sqlite3', got '%s'", cfg.SQL.Type)
	}
	if cfg.SQL.DB != "./toml.db" {
		t.Errorf("Expected SQL.DB to be './toml.db', got '%s'", cfg.SQL.DB)
	}
	// ...additional assertions...
}

func TestLoadInvalidConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to invalid config, but did not panic")
		}
	}()

	once = sync.Once{} // Reset the once variable
	cfg := NewConfig()
	cfg.Load("testdata/invalid_config.yaml", false)
}

func TestLoadInvalidEnv(t *testing.T) {
	// Set an invalid database type
	os.Setenv("EXT_SQL_TYPE", "invalid_db_type")
	defer func() {
		os.Unsetenv("EXT_SQL_TYPE")
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to invalid EXT_SQL_TYPE, but did not panic")
		}
	}()

	once = sync.Once{} // Reset the once variable
	cfg := NewConfig()
	cfg.Load("", true)
}

// CustomConfig represents test configuration structure
type CustomConfig struct {
	ServerName string   `mapstructure:"serverName"`
	Port       int      `mapstructure:"port"`
	Features   []string `mapstructure:"features"`
}

func TestCustomConfig(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     CustomConfig
	}{
		{
			name:     "load custom config from yaml",
			filepath: "testdata/config_with_custom.yaml",
			want: CustomConfig{
				ServerName: "test-server",
				Port:       8080,
				Features:   []string{"feature1", "feature2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset configuration and set expected custom config type
			cfg := setupTest()
			customCfg := &CustomConfig{}
			cfg.Custom = customCfg // Set the custom config type before loading

			// Load configuration
			cfg.Load(tt.filepath, false)

			// Verify custom config was loaded correctly
			got, ok := cfg.Custom.(*CustomConfig)
			if !ok {
				t.Fatal("Failed to get custom config as *CustomConfig")
			}

			// Compare all fields
			if got.ServerName != tt.want.ServerName {
				t.Errorf("ServerName = %v, want %v", got.ServerName, tt.want.ServerName)
			}
			if got.Port != tt.want.Port {
				t.Errorf("Port = %v, want %v", got.Port, tt.want.Port)
			}
			if !reflect.DeepEqual(got.Features, tt.want.Features) {
				t.Errorf("Features = %v, want %v", got.Features, tt.want.Features)
			}
		})
	}
}

func TestLoadServerConfig(t *testing.T) {
	// Set environment variables
	env := map[string]string{
		"EXT_SERVER_BINDADDR":       "127.0.0.1:8080",
		"EXT_SERVER_CORS":           "true",
		"EXT_SERVER_DOC":            "false",
		"EXT_SERVER_TRACE_ENABLED":  "true",
		"EXT_SERVER_TRACE_ENDPOINT": "http://localhost:4317",
	}

	// Set environment variables
	for k, v := range env {
		os.Setenv(k, v)
	}

	// Clean up environment variables after test
	defer func() {
		for k := range env {
			os.Unsetenv(k)
		}
	}()

	cfg := setupTest()
	cfg.Load("", true)

	// Validate server configuration
	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"Server.BindAddr", cfg.Server.BindAddr, "127.0.0.1:8080"},
		{"Server.Trace.Enabled", cfg.Server.Trace.Enabled, true},
		{"Server.Trace.Endpoint", cfg.Server.Trace.Endpoint, "http://localhost:4317"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLoadInvalidServerConfig(t *testing.T) {
	testCases := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{
			name: "invalid bind address",
			envVars: map[string]string{
				"EXT_SERVER_BINDADDR": "invalid:addr:8080",
			},
			expectedErr: "invalid address format",
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"EXT_SERVER_BINDADDR": "127.0.0.1:99999",
			},
			expectedErr: "invalid port number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			// Clean up environment variables
			defer func() {
				for k := range tc.envVars {
					os.Unsetenv(k)
				}
			}()

			defer func() {
				r := recover()
				if r == nil {
					t.Error("Expected panic, but got none")
					return
				}

				panicErr, ok := r.(error)
				if !ok {
					t.Errorf("Expected error panic, got %v", r)
					return
				}

				if !strings.Contains(panicErr.Error(), tc.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.expectedErr, panicErr.Error())
				}
			}()

			cfg := setupTest()
			cfg.Load("", true)
		})
	}
}
