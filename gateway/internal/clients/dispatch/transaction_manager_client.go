package dispatch

import (
	"context"
	"errors"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	transactionManagerClient struct {
		client grpc_gen.TransactionManagerClient
	}
)

var (
	ErrNotFound                            = errors.New("not found")
	ErrSubscriptionToDispatchAlreadyExists = errors.New("subscription already exists")
)

func NewTransactionManagerClient(conn grpc.ClientConnInterface) *transactionManagerClient {
	return &transactionManagerClient{
		client: grpc_gen.NewTransactionManagerClient(conn),
	}
}

func (c *transactionManagerClient) SubscribeForDispatch(ctx context.Context, email, dispatchID string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.NotFound {
		return errors.Join(ErrNotFound, err)
	}
	if status.Code(err) == codes.AlreadyExists {
		return ErrSubscriptionToDispatchAlreadyExists
	}

	return err
}

func (c *transactionManagerClient) UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.NotFound {
		return errors.Join(ErrNotFound, err)
	}

	return err
}
