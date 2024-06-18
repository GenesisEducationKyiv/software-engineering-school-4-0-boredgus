package controllers

import (
	"context"
	"fmt"
	"net/http"
	client_mocks "subscription-api/internal/mocks/clients"
	controllers_mocks "subscription-api/internal/mocks/controllers"
	"subscription-api/internal/services"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestSubscribeForDailyDispatch(t *testing.T) {
	type mocked struct {
		bindErr        error
		ctx            context.Context
		subscribeErr   error
		responseStatus int
		responseStr    string
	}
	dsClientMock := client_mocks.NewDispatchServiceClient(t)
	contextMock := controllers_mocks.NewContext(t)

	setup := func(m *mocked) func() {
		bindCall := contextMock.EXPECT().
			BindJSON(mock.Anything).Once().Return(m.bindErr)
		stringCall := contextMock.EXPECT().
			String(m.responseStatus, m.responseStr).NotBefore(bindCall).Maybe()
		contextCall := contextMock.EXPECT().
			Context().Once().NotBefore(bindCall).Return(m.ctx)
		subscribeCall := dsClientMock.EXPECT().
			SubscribeForDispatch(m.ctx, "", services.USD_UAH_DISPATCH_ID).Once().NotBefore(contextCall).
			Return(m.subscribeErr)
		statusCall := contextMock.EXPECT().
			Status(m.responseStatus).NotBefore(subscribeCall).Maybe()

		return func() {
			bindCall.Unset()
			stringCall.Unset()
			contextCall.Unset()
			subscribeCall.Unset()
			statusCall.Unset()
		}
	}

	someErr := fmt.Errorf("some err")
	tests := []struct {
		name   string
		mocked *mocked
	}{
		{
			name: "invalid data provided",
			mocked: &mocked{
				bindErr:        someErr,
				responseStr:    "invalid data provided",
				responseStatus: http.StatusBadRequest,
			},
		},
		{
			name: "user already subsccribed for this dispatch",
			mocked: &mocked{
				subscribeErr:   grpc.Errorf(codes.AlreadyExists, ""),
				responseStatus: http.StatusConflict,
			},
		},
		{
			name: "internal server error occured",
			mocked: &mocked{
				subscribeErr:   someErr,
				responseStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "success",
			mocked: &mocked{
				responseStatus: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked)
			defer cleanup()

			SubscribeForDailyDispatch(contextMock, dsClientMock)
		})
	}
}
