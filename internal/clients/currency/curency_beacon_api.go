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

const (
	CurrencyBeaconAPIBasePath string = "https://api.currencybeacon.com/v1/"
	CurrencyBeaconAPILabel    string = "Currency Beacon API https://currencybeacon.com/"
)

type currencyBeaconAPIClient struct {
	httpClient *clients.HTTPClient
	apiKey     string
	log        responseLogger
}

func NewCurrencyBeaconAPIClient(httpClient *clients.HTTPClient, apiKey string, logger config.Logger) *currencyBeaconAPIClient {
	return &currencyBeaconAPIClient{
		httpClient: httpClient,
		apiKey:     apiKey,
		log:        buildResponseLogger(logger, CurrencyBeaconAPILabel),
	}
}

// Gets latest exchange rates for specified currencies.
func (c *currencyBeaconAPIClient) Convert(
	ctx context.Context,
	baseCcy string,
	targetCcies []string,
) (map[string]float64, error) {
	resp, err := c.httpClient.Get(ctx, fmt.Sprintf("%slatest?api_key=%s&base=%s", CurrencyBeaconAPIBasePath, c.apiKey, baseCcy))
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	requestData := map[string]any{"base": baseCcy, "target": targetCcies}

	if resp.StatusCode != http.StatusOK {
		var parsedBody struct {
			Meta struct {
				ErrorType string `json:"error_type,omitempty"`
			} `json:"meta,omitempty"`
		}

		if err = utils.ParseJSON(resp.Body, &parsedBody); err != nil {
			c.log(resp.StatusCode, fmt.Errorf("failed to parse response: %w", err), requestData)

			return nil, ServiceIsUnaccessibleErr
		}

		err = fmt.Errorf(parsedBody.Meta.ErrorType)
		c.log(resp.StatusCode, err, requestData)

		return nil, err
	}

	var parsedBodyWithRates struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err = utils.ParseJSON(resp.Body, &parsedBodyWithRates); err != nil {
		// It's ok that on parse error we return UnsupportedCurrencyErr.
		// When unsupported currency is passed as base currency, API returns empty array as rates value.
		err := fmt.Errorf("%w: %s", UnsupportedCurrencyErr, baseCcy)
		c.log(resp.StatusCode, err, requestData)

		return nil, err
	}

	rates := make(map[string]float64)
	for _, currency := range targetCcies {
		ccy := strings.ToUpper(currency)
		if rate, ok := parsedBodyWithRates.Rates[ccy]; ok {
			rates[ccy] = rate
		}
	}

	requestData["rates"] = rates
	c.log(resp.StatusCode, nil, requestData)

	return rates, nil
}
