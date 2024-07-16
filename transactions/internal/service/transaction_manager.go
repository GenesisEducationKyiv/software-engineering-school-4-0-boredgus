package service

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/clients/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
)

type (
	CustomerService interface {
		CreateCustomer(ctx context.Context, email string) error
		CreateCustomerRevert(ctx context.Context, email string) error
	}
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatchID string) (*grpc_gen.Subscription, error)
		UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) (*grpc_gen.Subscription, error)
	}
	Broker interface {
		Publish(msg interface{})
	}
	transactionManager struct {
		transactionManagerURL string
		customerService       CustomerService
		dispatchService       DispatchService
		broker                Broker
		logger                config.Logger
	}
)

func NewTransactionManager(
	transactionManagerURL string,
	customerService CustomerService,
	dispatchService DispatchService,
	broker Broker,
	logger config.Logger,
) *transactionManager {
	return &transactionManager{
		transactionManagerURL: transactionManagerURL,
		customerService:       customerService,
		dispatchService:       dispatchService,
		broker:                broker,
		logger:                logger,
	}
}

func (m *transactionManager) SubscribeForDispatch(ctx context.Context, email, dispatchID string) error {
	if err := m.customerService.CreateCustomer(ctx, email); err != nil {
		m.logger.Errorf("failed to create a customer: %v", err)

		return err
	}

	subscription, err := m.dispatchService.SubscribeForDispatch(ctx, email, dispatchID)
	if err != nil {
		m.logger.Errorf("failed to subscribe for a dispatch: %v", err)

		if err = m.customerService.CreateCustomerRevert(ctx, email); err != nil {
			m.logger.Errorf("failed to revert creation of customer: %v", err)
		}

		return err
	}

	m.broker.Publish(subscription)

	return nil
}

func (m *transactionManager) UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) error {
	subscription, err := m.dispatchService.UnsubscribeFromDispatch(ctx, email, dispatchID)
	if err != nil {
		m.logger.Errorf("failed to subscribe for a dispatch: %v", err)

		return err
	}

	m.broker.Publish(subscription)

	return nil
}
