package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// 通用时间格式
const (
	TimeFormatDate     = "2006-01-02"
	TimeFormatDateTime = "2006-01-02 15:04:05"
)

type Config struct {
	Level     string
	Filepath  string
	MaxSizeMB int
	MaxAgeDay int
	Backups   int
	Compress  bool
}

func init() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()
}

func Init(cfg Config) {
	var zapLevel zapcore.Level

	// 日志等级解析
	switch cfg.Level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warning", "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// lumberjack 日志切割配置
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filepath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.Backups,
		MaxAge:     cfg.MaxAgeDay,
		Compress:   cfg.Compress,
	})

	// 基础编码配置 (提取公共部分)
	baseEencoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "Stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(TimeFormatDateTime))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 1. 控制台编码器（开启颜色）
	consoleEncoderConfig := baseEencoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConsole := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 2.文件编码器（json 格式, 无色）
	fileEncoderConfig := baseEencoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderJson := zapcore.NewJSONEncoder(fileEncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoderJson, writeSyncer, zapLevel),
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), zapLevel))

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()

	zap.ReplaceGlobals(logger)
}
