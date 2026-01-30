package log

import (
	"go.uber.org/zap"
)

// GetLogger 获取原生 zap logger
func GetLogger() *zap.Logger {
	return logger
}

// Sync 刷新缓冲区
func Sync() {
	_ = logger.Sync()
}

// ============================================================================
// 结构化日志 (Structured Logging) - 推荐高性能场景使用
// 用法: log.Info("user login", zap.String("username", "admin"))
// ============================================================================

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// ============================================================================
// 格式化日志 (Sugared Logging) - 方便使用，类似 fmt.Printf
// 用法: log.Infof("user %s login failed, count: %d", "admin", 3)
// ============================================================================

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}
