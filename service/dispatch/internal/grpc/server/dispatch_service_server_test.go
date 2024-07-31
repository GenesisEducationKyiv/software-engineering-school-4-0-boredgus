package server

import (
	"context"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"

	logger_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/logger"
	service_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_DispatchServiceServer_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc_gen.SubscribeForDispatchRequest
	}
	type mocked struct {
		subscription *entities.Subscription
		subscribeErr error
	}

	dsMock := service_mock.NewDispatchService(t)
	loggerMock := logger_mock.NewLogger()
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			SubscribeForDispatch(a.ctx, a.req.Email, a.req.DispatchId).
			Once().Return(m.subscription, m.subscribeErr)

		return func() {
			dsCall.Unset()
		}
	}

	arguments := &args{
		ctx: context.Background(),
		req: &grpc_gen.SubscribeForDispatchRequest{
			Email:      "email",
			DispatchId: "dispatch-id",
		},
	}

	subscription := entities.Subscription{}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc_gen.SubscribeForDispatchResponse
		expectedErr      error
	}{
		{
			name:         "failed: user already subscribed for this dispatch",
			args:         arguments,
			mockedValues: &mocked{subscribeErr: service.ErrUniqueViolation},
			expectedErr:  status.Error(codes.AlreadyExists, service.ErrUniqueViolation.Error()),
		},
		{
			name:         "failed: dispatch with such id does not exist",
			args:         arguments,
			mockedValues: &mocked{subscribeErr: service.ErrNotFound},
			expectedErr:  status.Error(codes.NotFound, service.ErrNotFound.Error()),
		},
		{
			name:         "failed: got unknown error from SubscribeForDispatch",
			args:         arguments,
			mockedValues: &mocked{subscribeErr: assert.AnError},
			expectedErr:  status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:         "success: subscribed for a dispatch",
			args:         arguments,
			mockedValues: &mocked{subscription: &subscription},
			expectedResponse: &grpc_gen.SubscribeForDispatchResponse{
				Subscription: subscriptionToProto(&subscription),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc_gen.UnimplementedDispatchServiceServer{},
			}
			actualResp, actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.expectedResponse, actualResp)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}

func Test_DispatchServiceServer_UnsubscribeFromDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc_gen.UnsubscribeFromDispatchRequest
	}
	type mocked struct {
		subscription   *entities.Subscription
		unsubscribeErr error
	}

	dsMock := service_mock.NewDispatchService(t)
	loggerMock := logger_mock.NewLogger()
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			UnsubscribeFromDispatch(a.ctx, a.req.Email, a.req.DispatchId).
			Once().Return(m.subscription, m.unsubscribeErr)

		return func() {
			dsCall.Unset()
		}
	}

	arguments := &args{
		ctx: context.Background(),
		req: &grpc_gen.UnsubscribeFromDispatchRequest{
			Email:      "email",
			DispatchId: "dispatch-id",
		},
	}

	subscription := entities.Subscription{
		Email:      arguments.req.Email,
		DispatchID: arguments.req.DispatchId,
		Status:     entities.SubscriptionStatusCancelled,
		SendAt:     time.Now(),
	}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc_gen.UnsubscribeFromDispatchResponse
		expectedErr      error
	}{
		{
			name:         "failed: dispatch with such id does not exist",
			args:         arguments,
			mockedValues: &mocked{unsubscribeErr: service.ErrNotFound},
			expectedErr:  status.Error(codes.NotFound, service.ErrNotFound.Error()),
		},
		{
			name:         "failed: got unknown error from SubscribeForDispatch",
			args:         arguments,
			mockedValues: &mocked{unsubscribeErr: assert.AnError},
			expectedErr:  status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:         "success: cancelled subscription",
			args:         arguments,
			mockedValues: &mocked{subscription: &subscription},
			expectedResponse: &grpc_gen.UnsubscribeFromDispatchResponse{
				Subscription: subscriptionToProto(&subscription),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc_gen.UnimplementedDispatchServiceServer{},
			}
			actualResp, actualErr := s.UnsubscribeFromDispatch(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.expectedResponse, actualResp)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
