package currency_client

import (
	"context"
	"errors"
	"fmt"
	"subscription-api/config"
	"subscription-api/internal/clients"
	"subscription-api/internal/services"
	"subscription-api/pkg/utils"
)

const (
	ExchangeRateAPILabel    string = "Exchange rate API https://exchangeratesapi.io/"
	ExchangeRateAPIBasePath string = "https://v6.exchangerate-api.com/v6"
)

type ExchangeRateAPIClient struct {
	apiKey     string
	httpClient *clients.HTTPClient
	logger     config.Logger
}

func NewExchangeRateAPIClient(httpClient *clients.HTTPClient, apiKey string, logger config.Logger) *ExchangeRateAPIClient {
	return &ExchangeRateAPIClient{
		apiKey:     apiKey,
		httpClient: httpClient,
		logger:     logger,
	}
}

var InvalidArgumentErr = errors.New("invalid argument")

var exchangeRateAPIErrors = map[string]string{
	"unsupported-code":  UnsupportedCurrencyErr.Error(),
	"malformed-request": "invalid structure of request",
	"invalid-key":       "invalid API key provided",
	"inactive-account":  "email address of account is not confirmed",
	"quota-reached":     "account has reached the the number of requests allowed by plan of subscription",
}

func (c *ExchangeRateAPIClient) logResponse(status int, err error, data map[string]any) {
	if err != nil {
		c.logger.Error(responseParams{
			Issuer:             ExchangeRateAPIBasePath,
			ResponseStatusCode: status,
			Error:              err.Error(),
			Data:               data,
		})

		return
	}
	c.logger.Info(responseParams{
		Issuer:             ExchangeRateAPIBasePath,
		ResponseStatusCode: status,
		Data:               data,
	})
}

// Gets latest exchange rates for specified currencies.
func (c *ExchangeRateAPIClient) Convert(
	ctx context.Context,
	baseCcy string,
	targetCcies []string,
) (map[string]float64, error) {
	resp, err := c.httpClient.Get(ctx, fmt.Sprintf("%s/%s/latest/%s", ExchangeRateAPIBasePath, c.apiKey, baseCcy))
	if err != nil {
		return nil, err
	}
	requestData := map[string]any{"base": baseCcy, "target": targetCcies}
	var parsedBody struct {
		ErrorType string             `json:"error-type,omitempty"`
		Rates     map[string]float64 `json:"conversion_rates,omitempty"`
	}
	if err = utils.ParseJSON(resp.Body, &parsedBody); err != nil {
		err = fmt.Errorf("failed to parse body: %w", err)
		c.logResponse(resp.StatusCode, err, requestData)

		return nil, err
	}
	if errorMsg, ok := exchangeRateAPIErrors[parsedBody.ErrorType]; ok {
		err = fmt.Errorf("%w: %v", services.FailedPreconditionErr, errorMsg)
		c.logResponse(resp.StatusCode, err, requestData)

		return nil, err
	}

	rates := make(map[string]float64)
	for _, currency := range targetCcies {
		rates[currency] = parsedBody.Rates[currency]
	}

	requestData["rates"] = rates
	c.logResponse(resp.StatusCode, nil, requestData)

	return rates, nil
}
