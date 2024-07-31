package controllers

import (
	"context"
	"net/http"
	"strconv"
)

type (
	CurrencyService interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}
	rateController struct {
		service CurrencyService
	}
)

func NewRateController(service CurrencyService) *rateController {
	return &rateController{
		service: service,
	}
}

func (c *rateController) GetExchangeRate(ctx Context) {
	from, to := "USD", "UAH"
	res, err := c.service.Convert(ctx.Context(), from, []string{to})
	if err != nil {
		ctx.Logger().Error(err)
		ctx.Status(http.StatusInternalServerError)

		return
	}
	ctx.String(http.StatusOK, strconv.FormatFloat(res[to], 'g', 7, 64))
}
