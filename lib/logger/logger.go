package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func Logger() *zap.Logger {
	return logger
}

func init() {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
}
