package entities

import (
	"time"
)

type CurrencyDispatchDetails struct {
	BaseCurrency     string
	TargetCurrencies []string
}

type Dispatch[T any] struct {
	ID                 string
	Label              string
	SendAt             time.Time
	TemplateName       string
	Details            T
	CountOfSubscribers int
}

type CurrencyDispatch Dispatch[CurrencyDispatchDetails]

func (d CurrencyDispatch) ToSubscription(email string, status SubscriptionStatus) *Subscription {
	return &Subscription{
		DispatchID:  d.ID,
		Email:       email,
		BaseCcy:     d.Details.BaseCurrency,
		TargetCcies: d.Details.TargetCurrencies,
		SendAt:      d.SendAt,
		Status:      status,
	}
}
