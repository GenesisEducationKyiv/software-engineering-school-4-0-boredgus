package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(...any)
	Errorf(format string, values ...any)
	Info(...any)
	Infof(format string, values ...any)
	Debug(...any)
	Debugf(format string, values ...any)
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.SugaredLogger.Infof(format, v...)
}

func InitLogger(mode Mode) *logger {
	config := zap.NewDevelopmentConfig()
	if mode == ProdMode {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return &logger{SugaredLogger: zap.Must(config.Build()).Sugar()}
}
