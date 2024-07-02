package controllers

import (
	"context"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch"
)

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
}

type subscribeParams struct {
	Email string `json:"email"`
}

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"

func SubscribeForDailyDispatch(ctx Context, ds DispatchService) {
	var params subscribeParams
	if err := ctx.BindJSON(&params); err != nil {
		ctx.String(http.StatusBadRequest, "invalid data provided")

		return
	}

	err := ds.SubscribeForDispatch(ctx.Context(), params.Email, USD_UAH_DISPATCH_ID)
	if err == dispatch.SubscriptionToDispatchAlreadyExistsErr {
		ctx.Status(http.StatusConflict)

		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
