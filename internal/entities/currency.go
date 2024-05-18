package entities

type Currency string

const (
	AmericanDollar   Currency = "USD"
	UkrainianHryvnia Currency = "UAH"
)

var SupportedCurrencies = []Currency{AmericanDollar, UkrainianHryvnia}

func (c Currency) IsSupported() bool {
	for _, cur := range SupportedCurrencies {
		if c == cur {
			return true
		}
	}
	return false
}
func FromString(data []string) []Currency {
	res := make([]Currency, len(data))
	for i, v := range data {
		res[i] = Currency(v)
	}
	return res
}
