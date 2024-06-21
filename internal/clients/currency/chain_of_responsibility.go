package currency_client

import "context"

type (
	CurrencyAPI interface {
		Convert(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error)
	}

	CurrencyAPIChain interface {
		CurrencyAPI
		SetNext(nextChain CurrencyAPIChain)
	}

	currencyAPIChain struct {
		api  CurrencyAPI
		next CurrencyAPIChain
	}
)

func NewCurrencyAPIChain(api CurrencyAPI) *currencyAPIChain {
	return &currencyAPIChain{api: api}
}

func (ch *currencyAPIChain) SetNext(chain CurrencyAPIChain) {
	ch.next = chain
}

func (ch *currencyAPIChain) Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error) {
	rates, err := ch.api.Convert(ctx, baseCcy, targetCcies)
	if err != nil {
		if ch.next == nil {
			return rates, err
		}

		return ch.next.Convert(ctx, baseCcy, targetCcies)
	}

	return rates, nil
}
