package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	e "subscription-api/internal/entities"
	"subscription-api/internal/services"
	"subscription-api/pkg/utils"
)

type HTTPClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type ExchangeRateAPIClient struct {
	basePath   string
	httpClient HTTPClient
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
	baseCcy e.Currency,
	targetCcies []e.Currency,
) (map[e.Currency]float64, error) {
	resp, err := c.httpClient.Get(ctx, fmt.Sprintf("%s/latest/%s", c.basePath, baseCcy))
	if err != nil {
		return nil, err
	}

	var result conversionResult
	if err = utils.ParseJSON(resp.Body, &result); err != nil {
		return nil, err
	}
	if result.ErrorType == InvalidArgumentErr.Error() {
		return nil, fmt.Errorf("%w: %v", services.InvalidArgumentErr, baseCcy)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", services.FailedPreconditionErr, targetCcies)
	}
	rates := make(map[e.Currency]float64)
	for _, currency := range targetCcies {
		rates[currency] = result.Rates[string(currency)]
	}

	return rates, nil
}
