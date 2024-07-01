package dispatch

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/clients/dispatch/gen"

	"google.golang.org/grpc"
)

type (
	dispatchServiceClient struct {
		client grpc_gen.DispatchServiceClient
	}

	Dispatch struct {
		Id                 string
		Label              string
		SendAt             string
		CountOfSubscribers int
	}
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
	// for _, dispatch := range dispatches {
	// 	convertedDispatches = append(convertedDispatches, DispatchData{
	// 		Id:                 dispatch.Id,
	// 		Label:              dispatch.Label,
	// 		SendAt:             dispatch.SendAt,
	// 		CountOfSubscribers: int(dispatch.CountOfSubscribers),
	// 	})
	// }

	return convertedDispatches
}
