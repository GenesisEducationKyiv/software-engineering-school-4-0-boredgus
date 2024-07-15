package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch"
)

type (
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
		UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) error
	}

	subscriptionRequest struct {
		Email string `json:"email"`
	}

	subscriptionController struct {
		service DispatchService
	}
)

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"

func NewSubscriptionController(service DispatchService) *subscriptionController {
	return &subscriptionController{
		service: service,
	}
}

func (c *subscriptionController) SubscribeForDailyDispatch(ctx Context) {
	var params subscriptionRequest
	if err := ctx.BindJSON(&params); err != nil {
		ctx.String(http.StatusBadRequest, "invalid data provided")

		return
	}

	err := c.service.SubscribeForDispatch(ctx.Context(), params.Email, USD_UAH_DISPATCH_ID)
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

func (c *subscriptionController) UnsubscribeFromDailyDispatch(ctx Context) {
	var params subscriptionRequest
	if err := ctx.BindJSON(&params); err != nil {
		ctx.String(http.StatusBadRequest, "invalid data provided")

		return
	}

	err := c.service.UnsubscribeFromDispatch(ctx.Context(), params.Email, USD_UAH_DISPATCH_ID)
	if errors.Is(err, dispatch.NotFoundErr) {
		ctx.Status(http.StatusNotFound)

		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
