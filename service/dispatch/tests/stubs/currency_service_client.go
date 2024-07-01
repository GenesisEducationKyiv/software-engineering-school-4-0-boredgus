package stubs

import (
	c "context"
)

type CurrencyServiceClient struct{}

func NewCurrencyServiceClient() *CurrencyServiceClient {
	return &CurrencyServiceClient{}
}

func (c *CurrencyServiceClient) Convert(ctx c.Context, baseCcy string, targetCcies []string) (map[string]float64, error) {
	return map[string]float64{
		"UAH": 40,
	}, nil
}
