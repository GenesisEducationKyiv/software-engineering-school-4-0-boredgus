package controllers

import (
	"context"
	"net/http"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, in *pb_ds.SubscribeForDispatchRequest, opts ...grpc.CallOption) (*pb_ds.SubscribeForDispatchResponse, error)
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

	_, err := ds.SubscribeForDispatch(ctx.Context(), &pb_ds.SubscribeForDispatchRequest{
		Email:      params.Email,
		DispatchId: USD_UAH_DISPATCH_ID,
	})
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
