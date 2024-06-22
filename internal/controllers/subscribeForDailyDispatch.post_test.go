package controllers

import (
	"context"
	"net/http"
	client_mocks "subscription-api/internal/mocks/clients"
	controllers_mocks "subscription-api/internal/mocks/controllers"
	"subscription-api/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_SubscribeForDailyDispatch_Controller(t *testing.T) {
	type mocked struct {
		expectedBindErr        error
		ctx                    context.Context
		expectedSubscribeErr   error
		expectedResponseStatus int
		expecteResponseStr     string
	}
	dsClientMock := client_mocks.NewDispatchServiceClient(t)
	contextMock := controllers_mocks.NewContext(t)

	setup := func(m *mocked) func() {
		bindCall := contextMock.EXPECT().
			BindJSON(mock.Anything).Once().Return(m.expectedBindErr)
		stringCall := contextMock.EXPECT().
			String(m.expectedResponseStatus, m.expecteResponseStr).NotBefore(bindCall).Maybe()
		contextCall := contextMock.EXPECT().
			Context().Once().NotBefore(bindCall).Return(m.ctx)
		subscribeCall := dsClientMock.EXPECT().
			SubscribeForDispatch(m.ctx, "", services.USD_UAH_DISPATCH_ID).Once().NotBefore(contextCall).
			Return(m.expectedSubscribeErr)
		statusCall := contextMock.EXPECT().
			Status(m.expectedResponseStatus).NotBefore(subscribeCall).Maybe()

		return func() {
			bindCall.Unset()
			stringCall.Unset()
			contextCall.Unset()
			subscribeCall.Unset()
			statusCall.Unset()
		}
	}

	tests := []struct {
		name   string
		mocked *mocked
	}{
		{
			name: "failed: got an error from BindJSON",
			mocked: &mocked{
				expectedBindErr:        assert.AnError,
				expecteResponseStr:     "invalid data provided",
				expectedResponseStatus: http.StatusBadRequest,
			},
		},
		{
			name: "failed: user already subsccribed for this dispatch",
			mocked: &mocked{
				expectedSubscribeErr:   status.Error(codes.AlreadyExists, ""),
				expectedResponseStatus: http.StatusConflict,
			},
		},
		{
			name: "failed: got unknown error from SubscribeForDispatch",
			mocked: &mocked{
				expectedSubscribeErr:   assert.AnError,
				expectedResponseStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "successfuly subscribed",
			mocked: &mocked{
				expectedResponseStatus: http.StatusOK,
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
