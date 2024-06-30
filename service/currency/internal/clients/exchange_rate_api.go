package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/pkg/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/pkg/utils"
)

const (
	ExchangeRateAPILabel    string = "Exchange rate API https://exchangeratesapi.io/"
	ExchangeRateAPIBasePath string = "https://v6.exchangerate-api.com/v6"
)

type ExchangeRateAPIClient struct {
	apiKey     string
	httpClient *HTTPClient
	log        responseLogger
}

func NewExchangeRateAPIClient(httpClient *HTTPClient, apiKey string, logger logger.Logger) *ExchangeRateAPIClient {
	return &ExchangeRateAPIClient{
		apiKey:     apiKey,
		httpClient: httpClient,
		log:        buildResponseLogger(logger, ExchangeRateAPILabel),
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
	defer resp.Body.Close()

	requestData := map[string]any{"base": baseCcy, "target": targetCcies}
	var parsedBody struct {
		ErrorType string             `json:"error-type,omitempty"`
		Rates     map[string]float64 `json:"conversion_rates,omitempty"`
	}
	if err = utils.ParseJSON(resp.Body, &parsedBody); err != nil {
		err = fmt.Errorf("failed to parse body: %w", err)
		c.log(resp.StatusCode, err, requestData)

		return nil, err
	}
	if errorMsg, ok := exchangeRateAPIErrors[parsedBody.ErrorType]; ok {
		err = fmt.Errorf("%w: %v", service.FailedPreconditionErr, errorMsg)
		c.log(resp.StatusCode, err, requestData)

		return nil, err
	}

	rates := make(map[string]float64)
	for _, currency := range targetCcies {
		rates[currency] = parsedBody.Rates[currency]
	}

	requestData["rates"] = rates
	c.log(resp.StatusCode, nil, requestData)

	return rates, nil
}
