package clients

import (
	"context"
	"errors"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dispatchServiceClient struct {
	client grpc_gen.DispatchServiceClient
}

func NewDispatchServiceClient(conn grpc.ClientConnInterface) *dispatchServiceClient {
	return &dispatchServiceClient{
		client: grpc_gen.NewDispatchServiceClient(conn),
	}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, email, dispatchID string) (*grpc_gen.Subscription, error) {
	resp, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.NotFound {
		return nil, errors.Join(service.ErrNotFound, err)
	}
	if status.Code(err) == codes.AlreadyExists {
		return nil, errors.Join(service.ErrAlreadyExists, err)
	}
	if err != nil {
		return nil, err
	}

	return resp.Subscription, err
}

func (c *dispatchServiceClient) UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) (*grpc_gen.Subscription, error) {
	resp, err := c.client.UnsubscribeFromDispatch(ctx, &grpc_gen.UnsubscribeFromDispatchRequest{
		Email:      email,
		DispatchId: dispatchID,
	})
	if status.Code(err) == codes.NotFound {
		return nil, errors.Join(service.ErrNotFound, err)
	}
	if err != nil {
		return nil, err
	}

	return resp.Subscription, err
}
