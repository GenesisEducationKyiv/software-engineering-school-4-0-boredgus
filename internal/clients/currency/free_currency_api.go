package currency_client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"subscription-api/config"
	"subscription-api/internal/clients"
	"subscription-api/pkg/utils"
)

type freeCurrencyAPIClient struct {
	httpClient *clients.HTTPClient
	log        responseLogger
}

const (
	FreeCurrencyAPILabel    = "Free Currency Exchange Rates API https://github.com/fawazahmed0/exchange-api"
	FreeCurrencyAPIBasePath = "http://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api"
)

func NewFreeCurrencyAPIClient(httpClient *clients.HTTPClient, logger config.Logger) *freeCurrencyAPIClient {
	return &freeCurrencyAPIClient{
		httpClient: httpClient,
		log:        buildResponseLogger(logger, ExchangeRateAPILabel),
	}
}

// Gets latest exchange rates for specified currencies.
func (c *freeCurrencyAPIClient) Convert(
	ctx context.Context,
	baseCcy string,
	targetCcies []string,
) (map[string]float64, error) {
	resp, err := c.httpClient.Get(ctx, fmt.Sprintf("%s@latest/v1/currencies/%s.json", FreeCurrencyAPIBasePath, baseCcy))
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	requestData := map[string]any{"base": baseCcy, "target": targetCcies}
	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("%w: %s", UnsupportedCurrencyErr, baseCcy)
		c.log(resp.StatusCode, err, requestData)

		return nil, fmt.Errorf("%w: %s", UnsupportedCurrencyErr, baseCcy)
	}
	if resp.StatusCode != http.StatusOK {
		c.log(resp.StatusCode, ServiceIsUnaccessibleErr, requestData)

		return nil, ServiceIsUnaccessibleErr
	}

	var result map[string]any
	if err = utils.ParseJSON(resp.Body, &result); err != nil {
		err = fmt.Errorf("failed to parse response: %w", err)
		c.log(resp.StatusCode, err, requestData)

		return nil, err
	}

	allRates := (result[strings.ToLower(baseCcy)]).(map[string]any)
	rates := make(map[string]float64)

	for _, currency := range targetCcies {
		ccy := strings.ToLower(currency)
		if rate, ok := allRates[ccy]; ok {
			rates[strings.ToUpper(currency)] = (rate).(float64)
		}
	}

	requestData["rates"] = rates
	c.log(resp.StatusCode, nil, requestData)

	return rates, nil
}
