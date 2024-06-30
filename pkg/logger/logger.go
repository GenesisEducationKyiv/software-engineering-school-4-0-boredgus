package logger

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
func (l *logger) Flush() {
	l.SugaredLogger.Sync()
}

type LoggerOption func(*zap.SugaredLogger)

func WithProcess(processName string) LoggerOption {
	return func(l *zap.SugaredLogger) {
		l = l.With("process", processName)
	}
}

type Mode string

const (
	DevMode  Mode = "dev"
	TestMode Mode = "test"
	ProdMode Mode = "prod"
)

func InitLogger(mode string, opts ...LoggerOption) *logger {
	config := zap.NewDevelopmentConfig()
	if mode == string(ProdMode) {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	l := zap.Must(config.Build()).Sugar()
	for _, opt := range opts {
		opt(l)
	}

	return &logger{SugaredLogger: l}
}
