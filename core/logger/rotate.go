package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/helper"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strconv"
	"time"
)

func LumberjackWriteSyncer() zapcore.WriteSyncer {
	logName := time.Now().Format("2006-01-02-15")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(globalInstance.BASE_CONFIG.ZapLog.Director, "system.logger."+logName),
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,

		Compress: false,
	}
	if globalInstance.BASE_CONFIG.ZapLog.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

func RotateWriteSyncer() zapcore.WriteSyncer {
	// 创建日志目录
	logDir := globalInstance.BASE_CONFIG.ZapLog.Director
	timePath := time.Now().Format("2006/01/02")
	logPath := path.Join(logDir, timePath)
	if ok, _ := helper.PathExists(logPath); !ok {
		fmt.Printf("日志目录%v不存在，创建日志目录\n", logPath)
		_ = os.MkdirAll(logPath, os.ModePerm)
	}
	logName := strconv.Itoa(time.Now().Hour()) + ".logger"
	rotateLogsLogger, _ := rotatelogs.New(
		path.Join(logPath, logName),
		rotatelogs.WithLinkName(globalInstance.BASE_CONFIG.ZapLog.LinkName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if globalInstance.BASE_CONFIG.ZapLog.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(rotateLogsLogger))
	}
	return zapcore.AddSync(rotateLogsLogger)
}
