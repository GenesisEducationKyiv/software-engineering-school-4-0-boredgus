package services

import (
	"context"
	e "subscription-api/internal/entities"
)

type ConvertCurrencyParams struct {
	Base   e.Currency
	Target []e.Currency
}

type CurrencyService interface {
	Convert(ctx context.Context, params ConvertCurrencyParams) (map[e.Currency]float64, error)
}

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatch string) error
	SendDispatch(ctx context.Context, dispatch string) error
	GetAllDispatches(ctx context.Context) ([]e.CurrencyDispatch, error)
}

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"
