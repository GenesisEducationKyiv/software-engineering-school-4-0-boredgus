package currency_service

import (
	"context"
	"subscription-api/internal/services"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(services.InvalidArgumentErr, "no target currencies provided")
	}

	return s.currencyAPIClient.Convert(ctx, params.Base, params.Target)
}
