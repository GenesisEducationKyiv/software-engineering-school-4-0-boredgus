package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"subscription-api/internal/services"
	"subscription-api/pkg/utils"
)

type ExchangeRateAPIClient struct {
	basePath   string
	httpClient *httpClient
}

func NewExchangeRateAPIClient(apiKey string) *ExchangeRateAPIClient {
	return &ExchangeRateAPIClient{
		basePath:   "https://v6.exchangerate-api.com/v6/" + apiKey,
		httpClient: NewHTTPClient(),
	}
}

type conversionResult struct {
	Result    string             `json:"result"`
	ErrorType string             `json:"error-type"`
	Rates     map[string]float64 `json:"conversion_rates"`
}

var InvalidArgumentErr = errors.New("invalid argument")

// Gets latest exchange rates for specifies currencies.
func (c *ExchangeRateAPIClient) Convert(
	ctx context.Context,
	baseCcy string,
	targetCcies []string,
) (map[string]float64, error) {
	resp, err := c.httpClient.Get(ctx, fmt.Sprintf("%s/latest/%s", c.basePath, baseCcy))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", services.FailedPreconditionErr, targetCcies)
	}

	var result conversionResult
	if err = utils.ParseJSON(resp.Body, &result); err != nil {
		return nil, err
	}
	if result.ErrorType == InvalidArgumentErr.Error() {
		return nil, fmt.Errorf("%w: %v", services.InvalidArgumentErr, baseCcy)
	}
	rates := make(map[string]float64)
	for _, currency := range targetCcies {
		rates[currency] = result.Rates[currency]
	}

	return rates, nil
}
