package cs

import (
	"context"
	"errors"
	"net/http"
	"subscription-api/internal/entities"
)

type ConvertParams struct {
	From entities.Currency
	To   []entities.Currency
}

type CurrencyService interface {
	Convert(ctx context.Context, params ConvertParams) (map[entities.Currency]float64, error)
}

var InvalidArgumentErr = errors.New("invalid argument")
var FailedPreconditionErr = errors.New("failed precondition")
var InvalidRequestErr = errors.New("invalid-request")

type currencyService struct {
	APIBasePath string
	client      http.Client
}

func NewCurrencyService(apiKey string) CurrencyService {
	return &currencyService{
		APIBasePath: "https://v6.exchangerate-api.com/v6/" + apiKey,
	}
}
