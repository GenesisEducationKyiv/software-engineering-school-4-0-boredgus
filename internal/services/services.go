package services

import (
	"context"
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

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatch string) error
	SendDispatch(ctx context.Context, dispatch string) error
	GetAllDispatches(ctx context.Context) ([]DispatchData, error)
}

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"
