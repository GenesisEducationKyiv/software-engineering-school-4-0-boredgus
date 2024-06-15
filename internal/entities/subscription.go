package entities

import (
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
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

func (d CurrencyDispatch) ToProto() *pb_ds.DispatchData {
	return &pb_ds.DispatchData{
		Id: d.Id,
		// BaseCurrency:       d.Details.BaseCurrency,
		// TargetCurrencies:   d.Details.TargetCurrencies,
		Label:              d.Label,
		SendAt:             d.SendAt.Format(time.TimeOnly),
		CountOfSubscribers: int64(d.CountOfSubscribers),
	}
}
