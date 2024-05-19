package config

import (
	"sync"

	"go.uber.org/zap"
)

type Logger interface {
	Error(...any)
	Errorf(format string, values ...any)
	Info(...any)
	Infof(format string, values ...any)
	Debug(...any)
	Debugf(format string, values ...any)
}

var logger zap.SugaredLogger
var onceLogger sync.Once

func InitLogger(mode Mode) *zap.SugaredLogger {
	l, er := zap.NewDevelopment()
	if mode == ProdMode {
		l, er = zap.NewProduction()
	}
	lg := zap.Must(l, er).Sugar()
	logger = *lg
	return lg
}

func Log() *zap.SugaredLogger {
	onceLogger.Do(func() {
		InitLogger(DevMode)
	})
	return &logger
}
