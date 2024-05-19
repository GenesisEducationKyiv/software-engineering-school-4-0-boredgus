package controllers

import (
	"fmt"
	"net/http"

	"subscription-api/internal/entities"
	pb_cs "subscription-api/pkg/grpc/currency_service"
)

func GetExchangeRate(ctx Context, cs pb_cs.CurrencyServiceClient) {
	from, to := string(entities.AmericanDollar), string(entities.UkrainianHryvnia)
	res, err := cs.Convert(ctx.Context(),
		&pb_cs.ConvertRequest{
			BaseCurrency:     from,
			TargetCurrencies: []string{to}})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.String(http.StatusOK, fmt.Sprint(res.Rates[to]))
}
