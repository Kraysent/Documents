package log

import (
	"fmt"

	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func InitLogger() error {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.Development = true

	logger, err := config.Build()
	if err != nil {
		return err
	}

	globalLogger = logger

	Info("Initialized logger")

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	fmt.Println(msg, fields)
	globalLogger.Fatal(msg, fields...)
}
