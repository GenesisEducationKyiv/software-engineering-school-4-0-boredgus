package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/pkg/entities"
	// "subscription-api/internal/entities"
	// "subscription-api/internal/services"
)

var (
	InvalidArgumentErr    = errors.New("invalid argument")
	FailedPreconditionErr = errors.New("failed precondition")
)

type ConvertCurrencyParams struct {
	Base   string
	Target []string
}

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

func (s *currencyService) Convert(ctx context.Context, params ConvertCurrencyParams) (map[string]float64, error) {
	if len(params.Target) == 0 {
		return nil, fmt.Errorf("%w: no target currencies provided", InvalidArgumentErr)
	}

	ccies, err := entities.MakeCurrencies(append(params.Target, params.Base))
	if err != nil {
		return nil, errors.Join(InvalidArgumentErr, err)
	}

	return s.currencyAPIClient.Convert(ctx, ccies[0], ccies[1:])
}
