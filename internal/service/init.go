package service

import (
	"context"
	"fmt"
	"gostart/internal/config"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ReadConfig() {
	config.ReadConfigByViper()
}

var DB *gorm.DB
var RedisDB *redis.Client
var Ctx context.Context = context.Background()

func ConnectDB() {
	// username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Configs.DB.User, config.Configs.DB.Password, config.Configs.DB.Host, config.Configs.DB.Port, config.Configs.DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Configs.Redis.Addr,
		Password: "", // no password
		DB:       0,  // use default DB
	})
	RedisDB = rdb
}

var ZapLogger *zap.Logger
var ZapSugar *zap.SugaredLogger

func ZapLogInit() {
	var core zapcore.Core
	env := config.Configs.Zap.Env

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)
	// 控制台等级
	levelDebug := zap.NewAtomicLevelAt(zap.DebugLevel)
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
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, levelDebug)
	if env == "prod" {
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

		// 文件输出
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Configs.Zap.OutputPaths,
			MaxSize:    100, // MB
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		})
		// 文件输出（错误日志）
		errWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Configs.Zap.ErrorOutputPaths,
			MaxSize:    100,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		})

		levelInfo := zap.NewAtomicLevelAt(zap.InfoLevel)
		levelError := zap.NewAtomicLevelAt(zap.ErrorLevel)

		fileCore := zapcore.NewCore(fileEncoder, fileWriter, levelInfo)
		fileErrCore := zapcore.NewCore(fileEncoder, errWriter, levelError)
		core = zapcore.NewTee(consoleCore, fileCore, fileErrCore)
	} else {
		core = zapcore.NewTee(consoleCore)
	}

	logger := zap.New(
		core,
		//zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		//zap.AddCallerSkip(1),
	)
	ZapLogger = logger
	ZapSugar = logger.Sugar()
	defer logger.Sync()
}
