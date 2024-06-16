package stubs

import (
	"subscription-api/internal/mailing"

	"github.com/stretchr/testify/mock"
)

type MailmanStub struct {
	mock.Mock
}

func NewMailmanStub() *MailmanStub {
	return &MailmanStub{}
}
func (m *MailmanStub) Send(email mailing.Email) error {
	m.Called(email)

	return nil
}
