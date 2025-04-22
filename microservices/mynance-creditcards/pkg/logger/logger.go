package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("[ERROR] [LOGGER] Failed to initialize logger")
	}
	Log = logger
}

func CloseLogger() {
	if Log != nil {
		Log.Sync()
	}
}
