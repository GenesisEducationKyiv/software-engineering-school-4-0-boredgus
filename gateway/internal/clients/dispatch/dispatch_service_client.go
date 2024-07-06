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
		client grpc_gen.DispatchServiceClient
	}

	DispatchData struct {
		Id                 string
		Label              string
		SendAt             string
		CountOfSubscribers int
	}
)

var (
	NotFoundErr                            = errors.New("not found")
	SubscriptionToDispatchAlreadyExistsErr = errors.New("subscription already exists")
)

func NewDispatchServiceClient(conn grpc.ClientConnInterface) *dispatchServiceClient {
	return &dispatchServiceClient{
		client: grpc_gen.NewDispatchServiceClient(conn),
	}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchId,
	})
	if status.Code(err) == codes.AlreadyExists {
		return SubscriptionToDispatchAlreadyExistsErr
	}

	return err
}

func (c *dispatchServiceClient) UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &grpc_gen.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchId,
	})
	if status.Code(err) == codes.NotFound {
		return errors.Join(NotFoundErr, err)
	}

	return err
}
