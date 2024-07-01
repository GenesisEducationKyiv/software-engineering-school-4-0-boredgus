package entities

import (
	"fmt"
	"slices"
	"strings"
)

const (
	AmericanDollar   string = "USD"
	UkrainianHryvnia string = "UAH"
)

var allSupportedCurrencies = []string{
	AmericanDollar,
	UkrainianHryvnia,
}

// MakeCurrencies validates provided data whether it is supportedd currency and returns .
func MakeCurrencies(currencies []string) ([]string, error) {
	ccies := make([]string, len(currencies))
	for i, ccy := range currencies {
		currency := strings.ToUpper(ccy)

		if !slices.Contains(allSupportedCurrencies, currency) {
			return nil, fmt.Errorf("%s is not supported currency", currency)
		}
		ccies[i] = currency
	}

	return ccies, nil
}
