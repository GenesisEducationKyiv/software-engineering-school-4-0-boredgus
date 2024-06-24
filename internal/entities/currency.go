package entities

import (
	"fmt"
	"slices"
	"strings"
)

type Currency string

const (
	AmericanDollar   Currency = "USD"
	UkrainianHryvnia Currency = "UAH"
)

// GetAllCurrencies returns all supported currencies
func GetAllCurrencies() []Currency {
	return []Currency{AmericanDollar, UkrainianHryvnia}
}

// ValidateCurrencies checks every element of passed array whether it is supported currency.
func ValidateCurrencies(currencies []string) error {
	allCurrencies := GetAllCurrencies()

	for _, ccy := range currencies {
		currency := Currency(strings.ToUpper(ccy))

		if !slices.Contains(allCurrencies, currency) {
			return fmt.Errorf("%s is not supported currency", currency)
		}
	}

	return nil
}
