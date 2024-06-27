package config_mocks

type Logger struct{}

func (_m *Logger) Debug(values ...interface{}) {}

func (_m *Logger) Debugf(format string, values ...interface{}) {}

func (_m *Logger) Error(values ...interface{}) {}

func (_m *Logger) Errorf(format string, values ...interface{}) {}

func (_m *Logger) Info(params ...interface{}) {}

func (_m *Logger) Infof(format string, values ...interface{}) {}

// NewLogger creates a mock instance for config.Logger interface.
func NewLogger() *Logger {
	return &Logger{}
}
