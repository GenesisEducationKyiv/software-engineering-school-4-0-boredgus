package stubs

import (
	// "subscription-api/internal/mailing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mailing"
	"github.com/stretchr/testify/mock"
)

type MailmanStub struct {
	mock.Mock
}

// Creates a mock for interface mailing.Mailman.
func NewMailmanStub() *MailmanStub {
	return &MailmanStub{}
}
func (m *MailmanStub) Send(email mailing.Email) error {
	args := m.Called(email)
	returnedErr := args.Get(0)

	if returnedErr != nil {
		return returnedErr.(error)
	}

	return nil
}
