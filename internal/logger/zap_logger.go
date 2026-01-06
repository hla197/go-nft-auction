package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 是 Logger 接口的具体实现
type ZapLogger struct {
	sugaredLogger *zap.SugaredLogger
	// 缓存 Writer，避免重复创建
	ioWriter io.Writer
	once     sync.Once
	base     *zap.Logger
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}
func (l *ZapLogger) Infoln(args ...interface{}) {
	l.sugaredLogger.Infoln(args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}
func (l *ZapLogger) Debugln(args ...interface{}) {
	l.sugaredLogger.Debugln(args...)
}
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}
func (l *ZapLogger) Errorln(args ...interface{}) {
	l.sugaredLogger.Errorln(args...)
}
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}
func (l *ZapLogger) Panicln(args ...interface{}) {
	l.sugaredLogger.Panicln(args...)
}

func (l *ZapLogger) Close() {
	if l.sugaredLogger != nil {
		_ = l.sugaredLogger.Sync() // stdout 下必须忽略错误
	}
}

func (l *ZapLogger) Init() {
	l.once.Do(func() {
		encoder := getEncoder()
		ws := getLogWriter()

		level := zap.NewAtomicLevelAt(zap.InfoLevel)

		core := zapcore.NewCore(encoder, ws, level)

		logger := zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)

		// 1️⃣ 接管标准库 log
		std := zap.NewStdLog(logger)
		log.SetOutput(std.Writer())
		log.SetFlags(0)

		l.base = logger
		l.sugaredLogger = logger.Sugar()
		l.ioWriter = std.Writer()
	})
}

// 将 Gin 框架的输出重定向到 Zap
func (l *ZapLogger) GetIoWriter() io.Writer {
	return l.ioWriter
}

// JSON 编码配置
func getEncoder() zapcore.Encoder {
	// 日志输出为控制台格式
	// 自定义控制台编码配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
	}

	return zapcore.NewConsoleEncoder(encoderConfig)

	// 日志输出为json格式
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// 日志写入配置（同时写入文件和控制台）
func getLogWriter() zapcore.WriteSyncer {
	// 配置 lumberjack 进行日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/gin.log", // 日志文件路径
		MaxSize:    10,               // 每个文件最大 10MB
		MaxBackups: 5,                // 最多保留 5 个备份
		MaxAge:     30,               // 文件最多保存 30 天
		Compress:   true,             // 是否压缩
	}

	// 同时输出到文件和控制台
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
