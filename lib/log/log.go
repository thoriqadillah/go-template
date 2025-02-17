package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func Logger() *zap.Logger {
	return logger
}

func init() {
	logger, _ = zap.NewProduction()
}
