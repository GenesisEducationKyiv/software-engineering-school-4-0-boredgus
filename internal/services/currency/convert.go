package services

import (
	"context"
	"fmt"
	"net/http"
	"subscription-api/internal/entities"
	"subscription-api/pkg/tools"
)

type conversionResult struct {
	Result    string             `json:"result"`
	ErrorType string             `json:"error-type"`
	Rates     map[string]float64 `json:"conversion_rates"`
}

func (e *currencyService) Convert(ctx context.Context, params ConvertParams) (map[entities.Currency]float64, error) {
	if len(params.To) == 0 {
		return nil, fmt.Errorf("%w: no target currencies provided", InvalidArgumentErr)
	}
	resp, err := http.Get(fmt.Sprintf("%s/latest/%s", e.APIBasePath, params.From))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", InvalidRequestErr, err)
	}
	fmt.Println(resp.StatusCode)
	var res conversionResult
	if err = tools.Parse(resp.Body, &res); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && res.ErrorType == InvalidArgumentErr.Error() {
		return nil, fmt.Errorf("%w: %v", InvalidArgumentErr, params.From)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", FailedPreconditionErr, res.ErrorType)
	}
	rates := make(map[entities.Currency]float64)
	for _, currency := range params.To {
		rates[currency] = res.Rates[string(currency)]
	}
	fmt.Println("SUCCESS", rates)
	return rates, nil
}
