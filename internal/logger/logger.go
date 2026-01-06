package logger

import (
	"io"
)

var Log Logger

// Logger 是一个通用日志接口，类似于 fmt 接口风格

type Logger interface {
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Close()
	// 将 Gin 框架的输出重定向到 Zap
	GetIoWriter() io.Writer
	// 初始化
	Init()
}

func InitLogger() {
	// 初始化日志系统，使用Zap作为底层日志实现
	Log = &ZapLogger{}
	Log.Init()
}
