package services

import (
	"context"
	dispatch_grpc "subscription-api/internal/services/dispatch/grpc"
)

type ConvertCurrencyParams struct {
	Base   string
	Target []string
}

type CurrencyService interface {
	Convert(ctx context.Context, params ConvertCurrencyParams) (map[string]float64, error)
}

type DispatchData struct {
	Id                 string
	Label              string
	SendAt             string
	CountOfSubscribers int
}

func (d DispatchData) ToProto() *dispatch_grpc.DispatchData {
	return &dispatch_grpc.DispatchData{
		Id:                 d.Id,
		Label:              d.Label,
		SendAt:             d.SendAt,
		CountOfSubscribers: int64(d.CountOfSubscribers),
	}
}

func DispatchDataFromProto(dispatches []*dispatch_grpc.DispatchData) []DispatchData {
	convertedDispatches := make([]DispatchData, 0, len(dispatches))
	for _, dispatch := range dispatches {
		convertedDispatches = append(convertedDispatches, DispatchData{
			Id:                 dispatch.Id,
			Label:              dispatch.Label,
			SendAt:             dispatch.SendAt,
			CountOfSubscribers: int(dispatch.CountOfSubscribers),
		})
	}

	return convertedDispatches
}

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatch string) error
	SendDispatch(ctx context.Context, dispatch string) error
	GetAllDispatches(ctx context.Context) ([]DispatchData, error)
}
