package grpc

import (
	"context"
	"subscription-api/internal/services"
	dispatch_client "subscription-api/internal/services/dispatch/grpc"

	"google.golang.org/grpc"
)

type DispatchServiceClient interface {
	SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
	SendDispatch(ctx context.Context, dispatchId string) error
	GetAllDispatches(ctx context.Context) ([]services.DispatchData, error)
}

type dispatchServiceClient struct {
	client dispatch_client.DispatchServiceClient
}

func NewDispatchServiceClient(conn grpc.ClientConnInterface) *dispatchServiceClient {
	return &dispatchServiceClient{}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &dispatch_client.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchId,
	})

	return err
}

func (c *dispatchServiceClient) SendDispatch(ctx context.Context, dispatchId string) error {
	_, err := c.client.SendDispatch(ctx, &dispatch_client.SendDispatchRequest{
		DispatchId: dispatchId,
	})

	return err
}

func (c *dispatchServiceClient) GetAllDispatches(ctx context.Context) ([]services.DispatchData, error) {
	resp, err := c.client.GetAllDispatches(ctx, &dispatch_client.GetAllDispatchesRequest{})
	if err != nil {
		return nil, err
	}

	return services.DispatchDataFromProto(resp.Dispatches), nil
}
