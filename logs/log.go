package logs

import (
	"github.com/zhiyunliu/golibs/xlog"
)

var (
	logger xlog.Logger
)

func init() {
	logger = xlog.New()
}

func Info(args ...interface{}) {
	logger.Log(xlog.LevelInfo, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Logf(xlog.LevelInfo, format, args...)
}

func Error(args ...interface{}) {
	logger.Log(xlog.LevelError, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Logf(xlog.LevelError, format, args...)
}
