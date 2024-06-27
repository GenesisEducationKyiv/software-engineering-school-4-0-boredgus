package dispatch_server

import (
	"context"
	config_mocks "subscription-api/internal/mocks/config"
	services_mocks "subscription-api/internal/mocks/services"
	"subscription-api/internal/services"
	grpc "subscription-api/internal/services/dispatch/server/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_DispatchServiceServer_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc.SubscribeForDispatchRequest
	}
	type mocked struct {
		expectedSubscribeErr error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger()
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			SubscribeForDispatch(a.ctx, a.req.Email, a.req.DispatchId).
			Once().Return(m.expectedSubscribeErr)

		return func() {
			dsCall.Unset()
		}
	}

	arguments := &args{
		ctx: context.Background(),
		req: &grpc.SubscribeForDispatchRequest{
			Email:      "email",
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name         string
		args         *args
		mockedValues *mocked
		want         *grpc.SubscribeForDispatchResponse
		wantErr      error
	}{
		{
			name:         "failed: user already subscribed for this dispatch",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: services.UniqueViolationErr},
			wantErr:      status.Error(codes.AlreadyExists, services.UniqueViolationErr.Error()),
		},
		{
			name:         "failed: dispatch with such id does not exist",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: services.NotFoundErr},
			wantErr:      status.Error(codes.NotFound, services.NotFoundErr.Error()),
		},
		{
			name:         "failed: got unknown error from SubscribeForDispatch",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: assert.AnError},
			wantErr:      status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:         "successfuly subscribed for a dispatch",
			args:         arguments,
			mockedValues: &mocked{},
			want:         &grpc.SubscribeForDispatchResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc.UnimplementedDispatchServiceServer{},
			}
			got, err := s.SubscribeForDispatch(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_DispatchServiceServer_SendDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc.SendDispatchRequest
	}
	type mocked struct {
		expectedSendErr error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger()

	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			SendDispatch(a.ctx, a.req.DispatchId).
			Once().Return(m.expectedSendErr)

		return func() {
			dsCall.Unset()
		}
	}

	arguments := &args{
		ctx: context.Background(),
		req: &grpc.SendDispatchRequest{
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc.SendDispatchResponse
		expectedErr      error
	}{
		{
			name:         "failed: dispatch withsuch id does not exist",
			args:         arguments,
			mockedValues: &mocked{expectedSendErr: services.NotFoundErr},
			expectedErr:  status.Error(codes.NotFound, services.NotFoundErr.Error()),
		},
		{
			name:         "failed: got unknown error from SendDispatch",
			args:         arguments,
			mockedValues: &mocked{expectedSendErr: assert.AnError},
			expectedErr:  status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:             "dispatch was successfuly sent",
			args:             arguments,
			mockedValues:     &mocked{},
			expectedResponse: &grpc.SendDispatchResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc.UnimplementedDispatchServiceServer{},
			}
			actualResponse, actualErr := s.SendDispatch(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.expectedResponse, actualResponse)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}

func Test_DispatchServiceServer_GetAllDispatches(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc.GetAllDispatchesRequest
	}
	type mocked struct {
		expectedDispatches []services.DispatchData
		expectedGetErr     error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger()
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			GetAllDispatches(a.ctx).
			Once().Return(m.expectedDispatches, m.expectedGetErr)

		return func() {
			dsCall.Unset()
		}
	}

	arguments := &args{
		ctx: context.Background(),
		req: &grpc.GetAllDispatchesRequest{},
	}
	dispatches := []services.DispatchData{{
		Id:                 "id",
		Label:              "label",
		CountOfSubscribers: 2,
	}}
	dispatchProtos := make([]*grpc.DispatchData, 0, len(dispatches))
	for _, d := range dispatches {
		dispatchProtos = append(dispatchProtos, ToProtoDispatch(d))
	}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc.GetAllDispatchesResponse
		expectedErr      error
	}{
		{
			name:         "failed: got unknown error from GetAllDispatches",
			args:         arguments,
			mockedValues: &mocked{expectedGetErr: assert.AnError},
			expectedErr:  status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:         "success: there is no dispatches",
			args:         arguments,
			mockedValues: &mocked{},
			expectedResponse: &grpc.GetAllDispatchesResponse{
				Dispatches: []*grpc.DispatchData{},
			},
		},
		{
			name:         "success: got all dispatches",
			args:         arguments,
			mockedValues: &mocked{expectedDispatches: dispatches},
			expectedResponse: &grpc.GetAllDispatchesResponse{
				Dispatches: dispatchProtos,
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
				UnimplementedDispatchServiceServer: grpc.UnimplementedDispatchServiceServer{},
			}
			actualResponse, actualErr := s.GetAllDispatches(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.expectedResponse, actualResponse)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
