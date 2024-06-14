package cs

import (
	"context"
	"fmt"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type ConvertParams struct {
	From entities.Currency
	To   []entities.Currency
}

type CurrencyService interface {
	Convert(ctx context.Context, params ConvertCurrencyParams) (map[entities.Currency]float64, error)
}

type ConvertCurrencyParams struct {
	Base   entities.Currency
	Target []entities.Currency
}
type CurrencyAPIClient interface {
	Convert(ctx context.Context, baseCcy entities.Currency, targetCcies []entities.Currency) (map[entities.Currency]float64, error)
}
type currencyService struct {
	currencyAPIClient CurrencyAPIClient
}

func NewCurrencyService(client CurrencyAPIClient) CurrencyService {
	return &currencyService{
		currencyAPIClient: client,
	}
}

func (s *currencyService) Convert(ctx context.Context, params ConvertCurrencyParams) (map[entities.Currency]float64, error) {
	if len(params.Target) == 0 {
		return nil, fmt.Errorf("%w: no target currencies provided", services.InvalidArgumentErr)
	}

	return s.currencyAPIClient.Convert(ctx, params.Base, params.Target)
}
