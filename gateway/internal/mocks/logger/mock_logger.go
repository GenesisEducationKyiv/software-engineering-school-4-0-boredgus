package logger_mock

type loggerMock struct{}

func NewLoggerMock() *loggerMock {
	return &loggerMock{}
}

func (m *loggerMock) Error(...any)                        {}
func (m *loggerMock) Errorf(format string, values ...any) {}
func (m *loggerMock) Info(...any)                         {}
func (m *loggerMock) Infof(format string, values ...any)  {}
func (m *loggerMock) Debug(...any)                        {}
func (m *loggerMock) Debugf(format string, values ...any) {}
