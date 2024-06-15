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

type LoggerOption func(*zap.SugaredLogger)

func WithProcess(processName string) LoggerOption {
	return func(l *zap.SugaredLogger) {
		l = l.With("process", processName)
	}
}

func InitLogger(mode Mode, opts ...LoggerOption) *logger {
	config := zap.NewDevelopmentConfig()
	if mode == ProdMode {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	l := zap.Must(config.Build()).Sugar()
	for _, opt := range opts {
		opt(l)
	}

	return &logger{SugaredLogger: l}
}
