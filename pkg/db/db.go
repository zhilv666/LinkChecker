package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// Config 定义通用数据配置
type Config struct {
	Type          string // sqlite3, mysql, postgres
	DSN           string // Data Source Name (file path or connection string)
	MaxIdleConns  int
	MaxOpenConns  int
	MaxLifetime   time.Duration
	TablePrefix   string
	SingularTable bool
	Debug         bool
}

// New 初始化 GORM 数据库连接
// 返回 db, cleanup 函数, 以及错误
func New(cfg *Config, zapLogger *zap.Logger) (*gorm.DB, func(), error) {
	if cfg.Type == "sqlite3" {
		dbDir := filepath.Dir(cfg.DSN)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
		}
	}

	gormLogger := zapgorm2.New(zapLogger)
	gormLogger.SetAsDefault()

	if cfg.Debug {
		gormLogger.LogLevel = logger.Info
		gormLogger.SlowThreshold = 200 * time.Millisecond
	} else {
		gormLogger.LogLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: cfg.SingularTable,
		},
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
	}

	var db *gorm.DB
	var err error
	switch cfg.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN), gormConfig)
	case "sqlite3":
		dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", cfg.DSN)
		db, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			zapLogger.Error("failed to close db connection", zap.Error(err))
		} else {
			zapLogger.Info("db connection closed")
		}
	}

	return db, cleanup, nil
}
