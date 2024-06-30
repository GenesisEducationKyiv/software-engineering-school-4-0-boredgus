package entities

import (
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
