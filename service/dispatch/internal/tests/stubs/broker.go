package stubs

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	"github.com/stretchr/testify/mock"
)

type BrokerStub struct {
	mock.Mock
}

func NewBrokerStub() *BrokerStub {
	return &BrokerStub{}
}

func (b *BrokerStub) CreateSubscription(sub deps.Subscription) {
	b.Called(sub)
}
