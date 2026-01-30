package configs

import (
	"path/filepath"
)

// Database 数据库配置
type Database struct {
	Type        string `json:"type" yaml:"type"`
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	Name        string `json:"name" yaml:"name"`
	DBFile      string `json:"db_file" yaml:"dbFile"`
	TablePrefix string `json:"table_prefix" yaml:"tablePrefix"`
	SSLMode     string `json:"ssl_mode" yaml:"sslMode"`
	DSN         string `json:"dsn" yaml:"dsn"`
}

// Cors 跨域配置
type Cors struct {
	AllowOrigins []string `json:"allow_origins" yaml:"allowOrigins"`
	AllowMethods []string `json:"allow_methods" yaml:"allowMethods"`
	AllowHeaders []string `json:"allow_headers" yaml:"allowHeaders"`
}

// Log 日志配置
type Log struct {
	Level     string `json:"level" yaml:"level"`
	Filepath  string `json:"filepath" yaml:"filepath"`
	MaxSizeMB int    `json:"max_size_mb" yaml:"maxSizeMb"`
	MaxAgeDay int    `json:"max_age_day" yaml:"maxAgeDay"`
	Backups   int    `json:"backups" yaml:"backups"`
	Compress  bool   `json:"compress" yaml:"compress"`
}

// Config 总配置
type Config struct {
	Database Database `json:"database" yaml:"database"`
	Cors     Cors     `json:"cors" yaml:"cors"`
	Log      Log      `json:"log" yaml:"log"`
}

// DefaultConfig 生成默认配置
func DefaultConfig() *Config {
	rootDir := "."
	dataDir := filepath.Join(rootDir, "data")
	logsDir := filepath.Join(rootDir, "logs")

	dbPath := filepath.Join(dataDir, "data.db")
	logPath := filepath.Join(logsDir, "app.log")

	return &Config{
		Database: Database{
			Type:        "sqlite3",
			Port:        0,
			TablePrefix: "lc_",
			DBFile:      dbPath,
		},
		Cors: Cors{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		},
		Log: Log{
			Level:     "debug",
			Filepath:  logPath,
			MaxSizeMB: 10,
			MaxAgeDay: 7,
			Backups:   3,
			Compress:  true,
		},
	}
}
