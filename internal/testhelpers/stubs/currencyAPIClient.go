package stubs

import (
	c "context"
	e "subscription-api/internal/entities"
)

type CurrencyAPIClient struct{}

func NewCurrencyAPIClient() *CurrencyAPIClient {
	return &CurrencyAPIClient{}
}

const DefaultUSDtoUAHRate = 40

func (c *CurrencyAPIClient) Convert(ctx c.Context, baseCcy e.Currency, targetCcies []e.Currency) (map[e.Currency]float64, error) {
	return map[e.Currency]float64{
		e.UkrainianHryvnia: DefaultUSDtoUAHRate,
	}, nil
}
