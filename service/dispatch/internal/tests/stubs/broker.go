package stubs

import (
	"github.com/stretchr/testify/mock"
)

type BrokerStub struct {
	mock.Mock
}

func NewBrokerStub() *BrokerStub {
	return &BrokerStub{}
}

func (b *BrokerStub) Publish(sub any) {
	b.Called(sub)
}
