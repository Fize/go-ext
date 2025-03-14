package config

import "testing"

func TestNewSQLConfig(t *testing.T) {
	tests := []struct {
		name    string
		opts    []SQLConfigOption
		wantErr bool
	}{
		{
			name:    "default config",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "valid mysql config",
			opts: []SQLConfigOption{
				WithType("mysql"),
				WithHost("localhost:3306"),
				WithUser("root"),
				WithPassword("password"),
				WithDB("testdb"),
			},
			wantErr: false,
		},
		{
			name: "invalid db type",
			opts: []SQLConfigOption{
				WithType("postgres"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewSQLConfig(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSQLConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && cfg == nil {
				t.Error("NewSQLConfig() returned nil config without error")
			}
		})
	}
}
