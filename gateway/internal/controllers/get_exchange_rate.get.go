package controllers

import (
	"context"
	"net/http"
	"strconv"
)

type CurrencyService interface {
	Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
}

func GetExchangeRate(ctx Context, cs CurrencyService) {
	from, to := "USD", "UAH"
	res, err := cs.Convert(ctx.Context(), from, []string{to})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}
	ctx.String(http.StatusOK, strconv.FormatFloat(res[to], 'g', 7, 64))
}
