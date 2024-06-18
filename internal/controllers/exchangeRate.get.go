package controllers

import (
	"context"
	"fmt"
	"net/http"

	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type CurrencyService interface {
	Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error)
}

func GetExchangeRate(ctx Context, cs CurrencyService) {
	from, to := string(entities.AmericanDollar), string(entities.UkrainianHryvnia)
	res, err := cs.Convert(ctx.Context(),
		services.ConvertCurrencyParams{
			Base:   from,
			Target: []string{to},
		})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}
	ctx.String(http.StatusOK, fmt.Sprint(res[to]))
}
