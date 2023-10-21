package log

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func InitLogger(stdoutPath, stderrPath string) error {
	defaultLogger := zap.L()
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.Development = true
	config.OutputPaths = []string{"stdout"}
	if stdoutPath != "" {
		if _, err := os.Stat(stdoutPath); !errors.Is(err, &fs.PathError{}) {
			defaultLogger.Warn(fmt.Sprintf("path '%s' provided to stdout does not exist", stdoutPath))
		} else {
			config.OutputPaths = append(config.OutputPaths, stdoutPath)
		}
	}
	config.ErrorOutputPaths = []string{"stderr"}
	if stderrPath != "" {
		if _, err := os.Stat(stderrPath); !errors.Is(err, &fs.PathError{}) {
			defaultLogger.Warn(fmt.Sprintf("path '%s' provided to stderr does not exist", stderrPath))
		} else {
			config.ErrorOutputPaths = append(config.ErrorOutputPaths, stderrPath)
		}
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	globalLogger = logger

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
	globalLogger.Fatal(msg, fields...)
}
