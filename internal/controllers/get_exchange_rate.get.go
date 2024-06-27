package controllers

import (
	"context"
	"net/http"
	"strconv"

	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type CurrencyService interface {
	Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error)
}

func GetExchangeRate(ctx Context, cs CurrencyService) {
	from, to := entities.AmericanDollar, entities.UkrainianHryvnia
	res, err := cs.Convert(ctx.Context(),
		services.ConvertCurrencyParams{
			Base:   from,
			Target: []string{to},
		})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}
	ctx.String(http.StatusOK, strconv.FormatFloat(res[to], 'g', 7, 64))
}
