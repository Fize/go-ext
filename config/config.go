package config

import (
	"sync"

	"github.com/kkyr/fig"
)

// default configuration
const (
	// default database type
	_defaultDBType = "sqlite3"
	// default database file
	_defaultDB = "./sqlite.db"
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

type dbType string

const (
	Mysql   dbType = "mysql"
	Sqlite3 dbType = "sqlite3"
)

// DBConfig
type DBConfig struct {
	// database type only support mysql and sqlite, default sqlite
	Type dbType `fig:"type"`
	// database host, include port such as 127.0.0.1:3306
	Host string `fig:"host"`
	// database user
	User string `fig:"user"`
	// database password
	Password string `fig:"password"`
	// database name
	DB           string `fig:"db"`
	MaxIdleConns int    `fig:"maxIdleConns"`
	MaxOpenConns int    `fig:"maxOpenConns"`
	// print raw sql
	SqlDebug bool `fig:"sqlDebug"`
}

// Email
type Email struct {
	Account  string `fig:"account"`
	SMTP     string `fig:"smtp"`
	Port     int    `fig:"port"`
	Password string `fig:"password"`
}

// Log
type Log struct {
	// log file path
	Filename string `fig:"filename"`
	// log file max size, unit MB
	MaxSize int `fig:"maxSize"`
	// log file max backups
	MaxBackups int `fig:"maxBackups"`
	// log file max age, unit day
	MaxAge int `fig:"maxAge"`
	// log file compress
	Compress bool `fig:"compress"`
	// log level
	Level string `fig:"level"`
	// log format
	Format string `fig:"format"`
	// log output
	Output string `fig:"output"`
}

// Config
type Config struct {
	DB    *DBConfig `fig:"db"`
	Email *Email    `fig:"email"`
	Log   *Log      `fig:"log"`
}

var (
	config *Config
	once   = sync.Once{}
)

func Load(dir, name string) {
	once.Do(func() {
		config = new(Config)
		err := fig.Load(config, fig.Dirs(dir), fig.File(name))
		if err != nil {
			panic(err)
		}
		if config.DB == nil {
			config.DB = new(DBConfig)
		}
		if config.DB.Type != Mysql && config.DB.Type != Sqlite3 {
			config.DB.Type = _defaultDBType
		}
		if config.DB.Type == Sqlite3 && config.DB.DB == "" {
			config.DB.DB = _defaultDB
		}
		if config.Log == nil {
			config.Log = new(Log)
		}
		if config.Log.Filename == "" {
			config.Log.Filename = _defaultLogPath
		}
		if config.Log.MaxSize == 0 {
			config.Log.MaxSize = _defaultLogMaxSize
		}
		if config.Log.MaxBackups == 0 {
			config.Log.MaxBackups = _defaultLogMaxBackups
		}
		if config.Log.MaxAge == 0 {
			config.Log.MaxAge = _defaultLogMaxAge
		}
		if config.Log.Level == "" {
			config.Log.Level = _defaultLogLevel
		}
		if config.Log.Format != "json" {
			config.Log.Format = _defaultLogFormat
		}
	})
}

func Read() *Config {
	return config
}
