package config

import (
	"fmt"
	"testing"
)

func TestNewLogConfig(t *testing.T) {
	tests := []struct {
		name    string
		opts    []LogConfigConfigOption
		wantErr bool
	}{
		{
			name:    "default config",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "valid config with string format",
			opts: []LogConfigConfigOption{
				WithLevel("debug"),
				WithFormat("string"),
			},
			wantErr: false,
		},
		{
			name: "valid config with json format",
			opts: []LogConfigConfigOption{
				WithLevel("info"),
				WithFormat("json"),
			},
			wantErr: false,
		},
		{
			name: "valid config with numeric level",
			opts: []LogConfigConfigOption{
				WithLevel("4"), // debug level
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			opts: []LogConfigConfigOption{
				WithFormat("yaml"),
			},
			wantErr: true,
		},
		{
			name: "invalid level string",
			opts: []LogConfigConfigOption{
				WithLevel("trace"),
			},
			wantErr: true,
		},
		{
			name: "invalid level number too small",
			opts: []LogConfigConfigOption{
				WithLevel("-1"),
			},
			wantErr: true,
		},
		{
			name: "invalid level number too large",
			opts: []LogConfigConfigOption{
				WithLevel("101"),
			},
			wantErr: true,
		},
		{
			name: "valid high number becomes debug",
			opts: []LogConfigConfigOption{
				WithLevel("99"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewLogConfig(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && cfg == nil {
				t.Error("NewLogConfig() returned nil config without error")
			}
		})
	}
}

func TestLogLevelConversion(t *testing.T) {
	tests := []struct {
		level    string
		wantNum  int
		wantBack string
		wantErr  bool
	}{
		{"debug", 4, "debug", false},
		{"info", 3, "info", false},
		{"warn", 2, "warn", false},
		{"error", 1, "error", false},
		{"fatal", 0, "fatal", false},
		{"invalid", -1, "", true},
		{"", -1, "", true},    // empty string
		{"-1", -1, "", true},  // negative number
		{"101", -1, "", true}, // too large number
		{"3.5", -1, "", true}, // invalid format
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			num, err := getLevelNum(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLevelNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if num != tt.wantNum {
				t.Errorf("getLevelNum() = %v, want %v", num, tt.wantNum)
			}

			if !tt.wantErr {
				str, err := getLevelString(num)
				if err != nil {
					t.Errorf("getLevelString() unexpected error = %v", err)
					return
				}
				if str != tt.wantBack {
					t.Errorf("getLevelString() = %v, want %v", str, tt.wantBack)
				}
			}
		})
	}
}

// Add new test for boundary conditions
func TestLogLevelBoundaries(t *testing.T) {
	tests := []struct {
		num       int
		wantLevel string
		wantErr   bool
	}{
		{-1, "", true},      // below minimum
		{0, "fatal", false}, // minimum valid
		{1, "error", false},
		{2, "warn", false},
		{3, "info", false},
		{4, "debug", false},
		{5, "debug", false}, // higher numbers become debug
		{100, "debug", false},
		{101, "", true}, // too high
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("level_%d", tt.num), func(t *testing.T) {
			got, err := getLevelString(tt.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLevelString(%d) error = %v, wantErr %v", tt.num, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantLevel {
				t.Errorf("getLevelString(%d) = %v, want %v", tt.num, got, tt.wantLevel)
			}
		})
	}
}
