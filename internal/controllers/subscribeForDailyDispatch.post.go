package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"subscription-api/config"
	"subscription-api/internal/services"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
)

type subscribeParams struct {
	Email string `json:"email"`
}

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"

func SubscribeForDailyDispatch(ctx Context, ds pb_ds.DispatchServiceClient) {
	var params subscribeParams
	if err := ctx.BindJSON(&params); err != nil {
		config.Log().Info("failed to bind subscribe params from json: " + err.Error())
		ctx.String(http.StatusBadRequest, "invalid data provided")
		return
	}

	_, err := ds.SubscribeFor(ctx.Context(), &pb_ds.SubscribeForRequest{
		Email:      params.Email,
		DispatchId: USD_UAH_DISPATCH_ID,
	})
	fmt.Println("controller SubscribeForDailyDispatch err: ", err, errors.Is(err, services.UniqueViolationErr))
	if err != nil && errors.Is(err, services.UniqueViolationErr) {
		ctx.String(http.StatusConflict, "such email already subscribed for USD to UAH dispatch")
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
