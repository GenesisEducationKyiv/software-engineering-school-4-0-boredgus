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

// ValidateCurrencies checks every element of passed array whether it is supported currency.
func ValidateCurrencies(currencies []string) error {
	for _, ccy := range currencies {
		currency := strings.ToUpper(ccy)

		if !slices.Contains(allSupportedCurrencies, currency) {
			return fmt.Errorf("%s is not supported currency", currency)
		}
	}

	return nil
}
