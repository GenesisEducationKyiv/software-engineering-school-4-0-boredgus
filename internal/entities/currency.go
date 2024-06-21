package entities

import "strings"

type Currency string

const (
	AmericanDollar   Currency = "USD"
	UkrainianHryvnia Currency = "UAH"
)

var SupportedCurrencies = map[Currency]struct{}{
	AmericanDollar:   {},
	UkrainianHryvnia: {},
}

func (c Currency) IsSupported() bool {
	_, ok := SupportedCurrencies[c]

	return ok
}

// CurrenciesFromString converts []string to []Currency.
func CurrenciesFromString(data []string) []Currency {
	res := make([]Currency, len(data))
	for i, v := range data {
		res[i] = Currency(strings.ToUpper(v))
	}

	return res
}
