package logging

import (
	"PlexWarp/internal/config"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	ServiceLogger *logrus.Logger
	AccessLogger  *logrus.Logger
)

// Init 初始化日志
func Init() {
	// 初始化服务日志
	ServiceLogger = logrus.New()
	ServiceLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 初始化访问日志
	AccessLogger = logrus.New()
	AccessLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 设置服务日志输出
	setupLogger(ServiceLogger, "service.log", config.Logger.ServiceLogger)

	// 设置访问日志输出
	setupLogger(AccessLogger, "access.log", config.Logger.AccessLogger)
}

// setupLogger 设置日志输出
func setupLogger(logger *logrus.Logger, filename string, setting config.BaseLoggerSetting) {
	var writers []io.Writer

	// 控制台输出
	if setting.Console {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if setting.File {
		logFile := filepath.Join(config.LogDir, filename)
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writers = append(writers, file)
		}
	}

	if len(writers) > 0 {
		logger.SetOutput(io.MultiWriter(writers...))
	}
}

// SetLevel 设置日志级别
func SetLevel(level logrus.Level) {
	if ServiceLogger != nil {
		ServiceLogger.SetLevel(level)
	}
	if AccessLogger != nil {
		AccessLogger.SetLevel(level)
	}
}

// 服务日志方法
func Debug(args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Debugf(format, args...)
	}
}

func Info(args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Infof(format, args...)
	}
}

func Warn(args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Warnf(format, args...)
	}
}

func Error(args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if ServiceLogger != nil {
		ServiceLogger.Errorf(format, args...)
	}
}

// 访问日志方法
func AccessInfo(args ...interface{}) {
	if AccessLogger != nil {
		AccessLogger.Info(args...)
	}
}

func AccessInfof(format string, args ...interface{}) {
	if AccessLogger != nil {
		AccessLogger.Infof(format, args...)
	}
}