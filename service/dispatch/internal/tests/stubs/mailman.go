package stubs

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	"github.com/stretchr/testify/mock"
)

type MailmanStub struct {
	mock.Mock
}

// Creates a mock for interface Mailman.
func NewMailmanStub() *MailmanStub {
	return &MailmanStub{}
}
func (m *MailmanStub) Send(email deps.Email) error {
	args := m.Called(email)
	returnedErr := args.Get(0)

	if returnedErr != nil {
		return returnedErr.(error)
	}

	return nil
}
