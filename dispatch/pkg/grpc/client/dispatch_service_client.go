package grpc

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"google.golang.org/grpc"
)

type DispatchServiceClient interface {
	SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
	SendDispatch(ctx context.Context, dispatchId string) error
	GetAllDispatches(ctx context.Context) ([]repo.DispatchData, error)
}

type (
	dispatchServiceClient struct {
		client grpc_gen.DispatchServiceClient
	}

	Dispatch repo.DispatchData
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

	return err
}

func (c *dispatchServiceClient) SendDispatch(ctx context.Context, dispatchId string) error {
	_, err := c.client.SendDispatch(ctx, &grpc_gen.SendDispatchRequest{
		DispatchId: dispatchId,
	})

	return err
}

func (c *dispatchServiceClient) GetAllDispatches(ctx context.Context) ([]Dispatch, error) {
	resp, err := c.client.GetAllDispatches(ctx, &grpc_gen.GetAllDispatchesRequest{})
	if err != nil {
		return nil, err
	}

	return protoToDispatchData(resp.Dispatches), nil
}

func protoToDispatchData(dispatches []*grpc_gen.DispatchData) []Dispatch {
	convertedDispatches := make([]Dispatch, 0, len(dispatches))
	for _, dispatch := range dispatches {
		convertedDispatches = append(convertedDispatches, Dispatch(repo.DispatchData{
			Id:                 dispatch.Id,
			Label:              dispatch.Label,
			SendAt:             dispatch.SendAt,
			CountOfSubscribers: int(dispatch.CountOfSubscribers),
		}))
	}

	return convertedDispatches
}
