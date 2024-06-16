package stubs

import (
	"subscription-api/internal/mailing"
)

type mailmanStub struct {
}

func NewMailmanStub() *mailmanStub {
	return &mailmanStub{}
}
func (m *mailmanStub) Send(email mailing.Email) error {
	return nil
}
