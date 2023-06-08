package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := "/tmp"
	name := "config_test.yaml"
	// 创建一个yaml格式的临时文件
	tmpfile, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// 写入一个配置文件，yaml格式
	configData := `db:
  type: sqlite3
  db: test.db
log:
  filename: test.log
  maxSize: 100
  maxBackups: 10
  maxAge: 30
  level: debug
  format: json`
	if _, err := tmpfile.Write([]byte(configData)); err != nil {
		t.Fatal(err)
	}

	// 加载配置文件
	err = Load(dir, name)
	if err != nil {
		t.Fatal(err)
	}

	// 检查配置是否正确加载
	if config.DB.Type != "sqlite3" {
		t.Errorf("DB.Type = %q; want %q", config.DB.Type, "sqlite3")
	}
	if config.DB.DB != "test.db" {
		t.Errorf("DB.DB = %q; want %q", config.DB.DB, "test.db")
	}
	if config.Log.Filename != "test.log" {
		t.Errorf("Log.Filename = %q; want %q", config.Log.Filename, "test.log")
	}
	if config.Log.MaxSize != 100 {
		t.Errorf("Log.MaxSize = %d; want %d", config.Log.MaxSize, 100)
	}
	if config.Log.MaxBackups != 10 {
		t.Errorf("Log.MaxBackups = %d; want %d", config.Log.MaxBackups, 10)
	}
	if config.Log.MaxAge != 30 {
		t.Errorf("Log.MaxAge = %d; want %d", config.Log.MaxAge, 30)
	}
	if config.Log.Level != "debug" {
		t.Errorf("Log.Level = %q; want %q", config.Log.Level, "debug")
	}
	if config.Log.Format != "json" {
		t.Errorf("Log.Format = %q; want %q", config.Log.Format, "json")
	}
}
