package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	_debugLevel = "debug"
	_infoLevel  = "info"
	_warnLevel  = "warn"
	_errorLevel = "error"
	_panicLevel = "panic"
	_fatalLevel = "fatal"
)

var (
	logger      *zap.SugaredLogger
	logLevelMap = map[string]zapcore.Level{
		_debugLevel: zapcore.DebugLevel,
		_infoLevel:  zapcore.InfoLevel,
		_warnLevel:  zapcore.WarnLevel,
		_errorLevel: zapcore.ErrorLevel,
		_panicLevel: zapcore.DPanicLevel, // 仅打印日志不触发panic
		_fatalLevel: zapcore.FatalLevel,
	}
	loglevel = zapcore.DebugLevel
)

func init() {
	// 初始化,定义日志文件名,日志级别
	core := zapcore.NewCore(getEncoder(), getWriteSyncer("kme"), logLevelMap[GetEnvDf("LOG_LEVEL", _debugLevel)])
	tempLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.DPanicLevel))
	logger = tempLogger.Sugar()
}

func SetLogLevel(l zapcore.Level) {
	loglevel = l
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder // 全路径编码器
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriteSyncer(filename string) zapcore.WriteSyncer {
	logPath := filepath.Join(abs(), "logs")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(err)
	}
	logPath = filepath.Join(logPath, strings.Replace(filename, ".log", "", -1)+".log")
	fmt.Println("log path: ", logPath)
	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		MaxAge:     7,       // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
}

func abs() string {
	root := filepath.Dir(GetEnvDf("LOG_PATH", os.Args[0]))
	p, _ := filepath.Abs(root)
	return p
}

func Sync() error {
	return logger.Sync()
}

func Debug(args ...interface{}) {
	if loglevel == zapcore.DebugLevel {
		logger.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if loglevel == zapcore.DebugLevel {
		logger.Debugf(format, args...)
	}
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func LogPanic(args ...interface{}) { // nolint
	logger.DPanic(args...)
}

func LogPanicf(format string, args ...interface{}) { // nolint
	logger.DPanicf(format, args...)
}

// Fatal 慎用,recover都捕获不到
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func GetEnvDf(key string, df ...string) string {
	value := os.Getenv(key)
	if value == "" && len(df) > 0 {
		return df[0]
	}
	return value
}
