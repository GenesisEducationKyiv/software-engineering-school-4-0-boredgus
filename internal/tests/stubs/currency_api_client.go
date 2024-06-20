package stubs

import (
	c "context"
)

type CurrencyAPIClient struct{}

func NewCurrencyAPIClient() *CurrencyAPIClient {
	return &CurrencyAPIClient{}
}

func (c *CurrencyAPIClient) Convert(ctx c.Context, baseCcy string, targetCcies []string) (map[string]float64, error) {
	return map[string]float64{
		"UAH": 40,
	}, nil
}
