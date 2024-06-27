package entities

import (
	grpc "subscription-api/internal/services/dispatch/grpc"
	"time"
)

type CurrencyDispatchDetails struct {
	BaseCurrency     string
	TargetCurrencies []string
}

type Dispatch[T any] struct {
	Id                 string
	Label              string
	SendAt             time.Time
	TemplateName       string
	Details            T
	CountOfSubscribers int
}

type CurrencyDispatch Dispatch[CurrencyDispatchDetails]

func (d CurrencyDispatch) ToProto() *grpc.DispatchData {
	return &grpc.DispatchData{
		Id:                 d.Id,
		Label:              d.Label,
		SendAt:             d.SendAt.Format(time.TimeOnly),
		CountOfSubscribers: int64(d.CountOfSubscribers),
	}
}
