package entities

import (
	"time"
)

type CurrencyDispatchDetails struct {
	BaseCurrency     string
	TargetCurrencies []string
}

type Dispatch[T any] struct {
	Id      string
	SendAt  time.Time
	Details T
}

type CurrencyDispatch struct {
	Dispatch[CurrencyDispatchDetails]
	CountOfSubscribers int
}
