package stubs

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

type dispatchStoreMock struct {
}

func NewDispatchStoreMock() *dispatchStoreMock {
	return &dispatchStoreMock{}
}

func (m *dispatchStoreMock) AddSubscription(ctx context.Context, sub entities.Subscription) error {
	return nil
}
func (m *dispatchStoreMock) CancelSubscription(ctx context.Context, sub entities.Subscription) error {
	return nil
}
