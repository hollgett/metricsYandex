package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Log = log
	return nil
}
