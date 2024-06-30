package deps

import "context"

type CurrencyServiceClient interface {
	Convert(ctx context.Context, base string, target []string) (map[string]float64, error)
}
