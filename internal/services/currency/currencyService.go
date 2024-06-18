package currency_service

import (
	"context"
	"fmt"
	ss "subscription-api/internal/services"
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

func (s *currencyService) Convert(ctx context.Context, params ss.ConvertCurrencyParams) (map[string]float64, error) {
	if len(params.Target) == 0 {
		return nil, fmt.Errorf("%w: no target currencies provided", ss.InvalidArgumentErr)
	}

	return s.currencyAPIClient.Convert(ctx, params.Base, params.Target)
}
