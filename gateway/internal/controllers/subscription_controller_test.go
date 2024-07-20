package controllers

import (
	"context"
	"net/http"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch"
	context_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/mocks/context"
	service_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/mocks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_SubscribeForDailyDispatch(t *testing.T) {
	type mocked struct {
		expectedBindErr        error
		ctx                    context.Context
		expectedSubscribeErr   error
		expectedResponseStatus int
		expecteResponseStr     string
	}
	dispatchServiceMock := service_mock.NewDispatchService(t)
	contextMock := context_mock.NewContext(t)

	setup := func(m *mocked) func() {
		bindCall := contextMock.EXPECT().
			BindJSON(mock.Anything).Once().Return(m.expectedBindErr)
		stringCall := contextMock.EXPECT().
			String(m.expectedResponseStatus, m.expecteResponseStr).NotBefore(bindCall).Maybe()
		contextCall := contextMock.EXPECT().
			Context().Once().NotBefore(bindCall).Return(m.ctx)
		subscribeCall := dispatchServiceMock.EXPECT().
			SubscribeForDispatch(m.ctx, "", USD_UAH_DISPATCH_ID).Once().NotBefore(contextCall).
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
				expectedSubscribeErr:   dispatch.SubscriptionToDispatchAlreadyExistsErr,
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

			NewSubscriptionController(dispatchServiceMock).SubscribeForDailyDispatch(contextMock)
		})
	}
}

func Test_Controller_UnsubscribeFromDailyDispatch(t *testing.T) {
	type mocked struct {
		expectedBindErr        error
		ctx                    context.Context
		expectedUnsubscribeErr error
		expectedResponseStatus int
		expecteResponseStr     string
	}
	dispatchServiceMock := service_mock.NewDispatchService(t)
	contextMock := context_mock.NewContext(t)

	setup := func(m *mocked) func() {
		bindCall := contextMock.EXPECT().
			BindJSON(mock.Anything).Once().Return(m.expectedBindErr)
		stringCall := contextMock.EXPECT().
			String(m.expectedResponseStatus, m.expecteResponseStr).NotBefore(bindCall).Maybe()
		contextCall := contextMock.EXPECT().
			Context().Once().NotBefore(bindCall).Return(m.ctx)
		subscribeCall := dispatchServiceMock.EXPECT().
			UnsubscribeFromDispatch(m.ctx, "", USD_UAH_DISPATCH_ID).Once().NotBefore(contextCall).
			Return(m.expectedUnsubscribeErr)
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
			name: "failed: subscription or dispatch is not found",
			mocked: &mocked{
				expectedUnsubscribeErr: dispatch.NotFoundErr,
				expectedResponseStatus: http.StatusNotFound,
			},
		},
		{
			name: "failed: got unknown error from SubscribeForDispatch",
			mocked: &mocked{
				expectedUnsubscribeErr: assert.AnError,
				expectedResponseStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "successfuly unsubscribed",
			mocked: &mocked{
				expectedResponseStatus: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked)
			defer cleanup()

			NewSubscriptionController(dispatchServiceMock).UnsubscribeFromDailyDispatch(contextMock)
		})
	}
}
