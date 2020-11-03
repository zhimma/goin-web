package core

import (
	"github.com/zhimma/goin-web/core/logger"
	globalInstance "github.com/zhimma/goin-web/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var level zapcore.Level

func Zap() (logger *zap.Logger) {
	switch globalInstance.BASE_CONFIG.ZapLog.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(buildZapCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(buildZapCore())
	}

	if globalInstance.BASE_CONFIG.ZapLog.ShowLine {
		logger.WithOptions(zap.AddCaller())
	}

	return logger
}

func buildZapCore() (core zapcore.Core) {
	return zapcore.NewCore(getEncoder(), logger.RotateWriteSyncer(), level)
}

func getEncoder() zapcore.Encoder {

	if globalInstance.BASE_CONFIG.ZapLog.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  globalInstance.BASE_CONFIG.ZapLog.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case globalInstance.BASE_CONFIG.ZapLog.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case globalInstance.BASE_CONFIG.ZapLog.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case globalInstance.BASE_CONFIG.ZapLog.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case globalInstance.BASE_CONFIG.ZapLog.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(globalInstance.BASE_CONFIG.ZapLog.Prefix + "2006/01/02 - 15:04:05"))
}
