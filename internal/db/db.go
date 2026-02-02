package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/model"
	"github.com/zhilv666/linkchecker/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// 全局 DB 实例
var db *gorm.DB

func Init(cfg configs.Database) error {
	// 0. 防御性检查：防止重复初始化
	if db != nil {
		return nil
	}

	// 1. 自动创建数据库文件的父目录 (关键优化)
	// 如果 cfg.DBFile 是 "data/data.db"，必须先保证 "data" 目录存在
	dbDir := filepath.Dir(cfg.DBFile)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
	}

	// 2. 配置日志
	zapL := log.GetLogger()
	gormLogger := zapgorm2.New(zapL)
	gormLogger.SetAsDefault()

	// 优化：根据 Zap 的级别动态调整 GORM 的级别
	// 如果 zapL 是 Debug 级别，GORM 也开启 Info (打印 SQL)
	if zapL.Core().Enabled(0) { // Check DebugLevel (Warning: zap levels are negative/complex, simple mapping below)
		// 简单策略：生产环境只记录慢查询和错误，开发环境记录所有 SQL
		// 这里我们设为 Warn，但在开发时你可能想改为 logger.Info
		// gormLogger.LogLevel = logger.Warn
		gormLogger.LogLevel = logger.Info
		gormLogger.SlowThreshold = 200 * time.Millisecond
	} else {
		gormLogger.LogLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true, // 建议开启：使用 'user' 而不是 'users'，通常更符合直觉
		},
		Logger: gormLogger,
		// 禁用默认事务，提升 30% 写入性能 (如果没有强事务需求)
		SkipDefaultTransaction: true,
	}

	var err error
	switch cfg.Type {
	case "sqlite3":
		// _pragma=busy_timeout(5000) 防止并发锁死
		// _pragma=journal_mode(WAL) 开启 WAL 模式提升并发读写性能
		dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", cfg.DBFile)
		db, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		log.Errorf("failed to connect database: %v", err)
		return err
	}

	// 3. 连接池配置
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100) // WAL 模式下可以设置较高
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 自动迁移
	if err := db.AutoMigrate(&model.LinkRecord{}); err != nil {
		log.Errorf("failed to migrate database: %v", err)
		return err
	}

	return nil
}

// GetDB 获取全局数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		// 这里可以用 log.Panic 提醒开发者必须先调用 Init
		log.Panic("database not initialized, please call data.Init() first")
	}
	return db
}

// Close 关闭数据库连接
func Close() {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("failed to get sql.DB: %v", err)
		return
	}

	log.Info("closing db connection...")
	if err := sqlDB.Close(); err != nil {
		log.Errorf("failed to close db: %v", err)
	}
}
