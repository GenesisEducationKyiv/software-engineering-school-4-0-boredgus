package stubs

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"

type notifierMock struct {
}

func NewNotifierMock() *notifierMock {
	return &notifierMock{}
}
func (m *notifierMock) Notify(service.Notification) error {
	return nil
}
