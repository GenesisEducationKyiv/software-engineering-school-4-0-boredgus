package controllers

import (
	"context"
	"net/http"
	"subscription-api/internal/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
}

type subscribeParams struct {
	Email string `json:"email"`
}

func SubscribeForDailyDispatch(ctx Context, ds DispatchService) {
	var params subscribeParams
	if err := ctx.BindJSON(&params); err != nil {
		ctx.String(http.StatusBadRequest, "invalid data provided")

		return
	}

	err := ds.SubscribeForDispatch(ctx.Context(), params.Email, services.USD_UAH_DISPATCH_ID)
	if status.Code(err) == codes.AlreadyExists {
		ctx.Status(http.StatusConflict)

		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
