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
	dispatchServiceClient struct {
		client grpc_gen.TransactionManagerClient
	}

	DispatchData struct {
		Id                 string
		Label              string
		SendAt             string
		CountOfSubscribers int
	}
)

var (
	ErrNotFound                            = errors.New("not found")
	ErrSubscriptionToDispatchAlreadyExists = errors.New("subscription already exists")
)

func NewDispatchServiceClient(conn grpc.ClientConnInterface) *dispatchServiceClient {
	return &dispatchServiceClient{
		client: grpc_gen.NewTransactionManagerClient(conn),
	}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, email, dispatchID string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.AlreadyExists {
		return ErrSubscriptionToDispatchAlreadyExists
	}

	return err
}

func (c *dispatchServiceClient) UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.NotFound {
		return errors.Join(ErrNotFound, err)
	}

	return err
}
