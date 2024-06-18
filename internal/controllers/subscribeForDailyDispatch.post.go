package controllers

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		ctx.Logger().Debug("failed to bind subscribe params from json: " + err.Error())
		ctx.String(http.StatusBadRequest, "invalid data provided")

		return
	}

	err := ds.SubscribeForDispatch(ctx.Context(), params.Email, USD_UAH_DISPATCH_ID)
	if status.Code(err) == codes.AlreadyExists {
		ctx.Status(http.StatusConflict)

		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Logger().Debugf("failed to subscribe user for dispatch: %v", err)

		return
	}

	ctx.Status(http.StatusOK)
}
