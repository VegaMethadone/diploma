package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	once   sync.Once
)

func InitFileLogger(logPath string) {
	once.Do(func() {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncodeLevel = zapcore.CapitalLevelEncoder

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}

		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(config),
				zapcore.AddSync(logFile),
				zap.InfoLevel,
			),
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(config),
				zapcore.AddSync(os.Stderr),
				zap.ErrorLevel,
			),
		)

		Logger = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	})
}

func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

func NewInfoMessage(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func NewDebugMessage(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func NewWarnMessage(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func NewErrMessage(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func NewDPanicMessage(msg string, fields ...zap.Field) {
	Logger.DPanic(msg, fields...)
}

func NewPanicMessage(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

func NewFatalMessage(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
