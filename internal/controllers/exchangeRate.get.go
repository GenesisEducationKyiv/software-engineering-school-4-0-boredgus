package controllers

import (
	"context"
	"fmt"
	"net/http"

	"subscription-api/internal/entities"
	pb_cs "subscription-api/pkg/grpc/currency_service"

	"google.golang.org/grpc"
)

type CurrencyService interface {
	Convert(ctx context.Context, in *pb_cs.ConvertRequest, opts ...grpc.CallOption) (*pb_cs.ConvertResponse, error)
}

func GetExchangeRate(ctx Context, cs CurrencyService) {
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
