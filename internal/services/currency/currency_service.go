package currency_service

import (
	"context"
	"fmt"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type CurrencyAPIClient interface {
	Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
}
type currencyService struct {
	currencyAPIClient CurrencyAPIClient
}

func NewCurrencyService(client CurrencyAPIClient) *currencyService {
	return &currencyService{
		currencyAPIClient: client,
	}
}

func (s *currencyService) Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error) {
	if len(params.Target) == 0 {
		return nil, fmt.Errorf("%w: no target currencies provided", services.InvalidArgumentErr)
	}
	for _, ccy := range entities.CurrenciesFromString(append(params.Target, params.Base)) {
		if !ccy.IsSupported() {
			return nil, fmt.Errorf("%w: currency %s is not supported", services.InvalidArgumentErr, ccy)
		}
	}

	return s.currencyAPIClient.Convert(ctx, params.Base, params.Target)
}
