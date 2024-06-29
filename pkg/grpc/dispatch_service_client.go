package grpc

import (
	"context"
	"subscription-api/internal/services"
	dispatch_grpc "subscription-api/internal/services/dispatch/server/grpc"

	"google.golang.org/grpc"
)

type DispatchServiceClient interface {
	SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
	SendDispatch(ctx context.Context, dispatchId string) error
	GetAllDispatches(ctx context.Context) ([]services.DispatchData, error)
}

type dispatchServiceClient struct {
	client dispatch_grpc.DispatchServiceClient
}

func NewDispatchServiceClient(conn grpc.ClientConnInterface) *dispatchServiceClient {
	return &dispatchServiceClient{}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := c.client.SubscribeForDispatch(ctx, &dispatch_grpc.SubscribeForDispatchRequest{
		Email:      email,
		DispatchId: dispatchId,
	})

	return err
}

func (c *dispatchServiceClient) SendDispatch(ctx context.Context, dispatchId string) error {
	_, err := c.client.SendDispatch(ctx, &dispatch_grpc.SendDispatchRequest{
		DispatchId: dispatchId,
	})

	return err
}

func (c *dispatchServiceClient) GetAllDispatches(ctx context.Context) ([]services.DispatchData, error) {
	resp, err := c.client.GetAllDispatches(ctx, &dispatch_grpc.GetAllDispatchesRequest{})
	if err != nil {
		return nil, err
	}

	return protoToDispatchData(resp.Dispatches), nil
}

func protoToDispatchData(dispatches []*dispatch_grpc.DispatchData) []services.DispatchData {
	convertedDispatches := make([]services.DispatchData, 0, len(dispatches))
	for _, dispatch := range dispatches {
		convertedDispatches = append(convertedDispatches, services.DispatchData{
			Id:                 dispatch.Id,
			Label:              dispatch.Label,
			SendAt:             dispatch.SendAt,
			CountOfSubscribers: int(dispatch.CountOfSubscribers),
		})
	}

	return convertedDispatches
}
