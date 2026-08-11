package test

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Zap *zap.Logger

func TestZapLogInit() {
	// 控制台编码器
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 文件编码器（JSON 结构化）
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)
	// 文件输出
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	})

	// 文件输出（错误日志）
	errWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/err.log",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	})

	// ========== 3. 创建多个 Core ==========
	// 使用相同的日志级别
	levelDebug := zap.NewAtomicLevelAt(zap.DebugLevel)
	levelInfo := zap.NewAtomicLevelAt(zap.InfoLevel)
	levelError := zap.NewAtomicLevelAt(zap.ErrorLevel)

	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, levelDebug)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, levelInfo)
	fileErrCore := zapcore.NewCore(fileEncoder, errWriter, levelError)

	// 合并多个 Core
	core := zapcore.NewTee(consoleCore, fileCore, fileErrCore)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)
	Zap = logger
	defer logger.Sync()
}
