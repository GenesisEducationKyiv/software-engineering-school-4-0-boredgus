package service

import (
	"context"
	"errors"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/gen"
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

var (
	ErrTransportProblem = errors.New("transport problem")
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotFound         = errors.New("not found")
)

const (
	MaxCountOfRetries int           = 3
	RetryInterval     time.Duration = 100 * time.Millisecond
)

func (m *transactionManager) SubscribeForDispatch(ctx context.Context, email, dispatchID string) error {
	var err error

	m.retryIfFails(func() error {
		e := m.customerService.CreateCustomer(ctx, email)
		err = e

		return e
	}, isTransportProblem)

	if err != nil && !errors.Is(err, ErrAlreadyExists) {
		m.logger.Errorf("failed to create a customer: %v", err)

		return err
	}

	var subscription *grpc_gen.Subscription
	m.retryIfFails(func() error {
		subscription, err = m.dispatchService.SubscribeForDispatch(ctx, email, dispatchID)

		return err
	}, isTransportProblem)

	if err != nil {
		m.logger.Errorf("failed to subscribe for a dispatch: %v", err)

		m.retryIfFails(func() error {
			e := m.customerService.CreateCustomerRevert(ctx, email)
			err = e

			return e
		}, isTransportProblem)
		if err != nil {
			m.logger.Errorf("failed to revert creation of customer: %v", err)
		}

		return err
	}

	m.broker.Publish(subscription)

	return nil
}

func (m *transactionManager) UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) error {
	var subscription *grpc_gen.Subscription
	var err error

	m.retryIfFails(func() error {
		subscription, err = m.dispatchService.UnsubscribeFromDispatch(ctx, email, dispatchID)

		return err
	}, isTransportProblem)

	if err != nil {
		m.logger.Errorf("failed to subscribe for a dispatch: %v", err)

		return err
	}

	m.broker.Publish(subscription)

	return nil
}

func (m *transactionManager) retryIfFails(action func() error, successCondition func(error) bool) {
	for i := 0; i < MaxCountOfRetries; i++ {
		if successCondition(action()) {
			return
		}

		time.Sleep(RetryInterval)
	}
}

func isTransportProblem(err error) bool {
	return !errors.Is(err, ErrTransportProblem)
}
