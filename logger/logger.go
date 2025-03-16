package logger

import "go.uber.org/zap"

var Loger *zap.Logger

func NewLoger() {
	var err error
	Loger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func NewInfoMessage(msg string, fields ...zap.Field) {
	Loger.Info(msg, fields...)
}

func NewDebugMessage(msg string, fields ...zap.Field) {
	Loger.Debug(msg, fields...)
}

func NewWarnMessage(msg string, fields ...zap.Field) {
	Loger.Warn(msg, fields...)
}

func NewErrMessage(msg string, fields ...zap.Field) {
	Loger.Error(msg, fields...)
}

func NewDPanicMessage(msg string, fields ...zap.Field) {
	Loger.DPanic(msg, fields...)
}

func NewPanicMessage(msg string, fields ...zap.Field) {
	Loger.Panic(msg, fields...)
}

func NewFatalMessage(msg string, fields ...zap.Field) {
	Loger.Fatal(msg, fields...)
}
