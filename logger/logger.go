package logger

import (
	"log"

	"go.uber.org/zap"
)

type Log struct {
	Logger *zap.Logger
}

func Newlogger() *Log {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	return &Log{
		Logger: logger,
	}
    
}
