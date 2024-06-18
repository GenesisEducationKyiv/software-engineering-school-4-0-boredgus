package stubs

import (
	c "context"
)

type CurrencyAPIClient struct{}

func NewCurrencyAPIClient() *CurrencyAPIClient {
	return &CurrencyAPIClient{}
}

const DefaultUSDtoUAHRate = 40

func (c *CurrencyAPIClient) Convert(ctx c.Context, baseCcy string, targetCcies []string) (map[string]float64, error) {
	return map[string]float64{
		"UAH": DefaultUSDtoUAHRate,
	}, nil
}
