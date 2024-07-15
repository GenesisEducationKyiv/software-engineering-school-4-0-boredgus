package stubs

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/stretchr/testify/mock"
)

type BrokerStub struct {
	mock.Mock
}

func NewBrokerStub() *BrokerStub {
	return &BrokerStub{}
}

func (b *BrokerStub) CreateSubscription(sub service.Subscription) {
	b.Called(sub)
}

func (b *BrokerStub) CancelSubscription(sub service.Subscription) {
	b.Called(sub)
}
