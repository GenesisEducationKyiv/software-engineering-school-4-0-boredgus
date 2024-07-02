package server

import (
	"context"
	"testing"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"

	logger_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/logger"
	service_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	service_err "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"
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
		expectedSubscribeErr error
	}

	dsMock := service_mock.NewDispatchService(t)
	loggerMock := logger_mock.NewLogger()
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
		req: &grpc_gen.SubscribeForDispatchRequest{
			Email:      "email",
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name         string
		args         *args
		mockedValues *mocked
		want         *grpc_gen.SubscribeForDispatchResponse
		wantErr      error
	}{
		{
			name:         "failed: user already subscribed for this dispatch",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: service_err.UniqueViolationErr},
			wantErr:      status.Error(codes.AlreadyExists, service_err.UniqueViolationErr.Error()),
		},
		{
			name:         "failed: dispatch with such id does not exist",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: service_err.NotFoundErr},
			wantErr:      status.Error(codes.NotFound, service_err.NotFoundErr.Error()),
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
			want:         &grpc_gen.SubscribeForDispatchResponse{},
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
		req *grpc_gen.SendDispatchRequest
	}
	type mocked struct {
		expectedSendErr error
	}

	dsMock := service_mock.NewDispatchService(t)
	loggerMock := logger_mock.NewLogger()

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
		req: &grpc_gen.SendDispatchRequest{
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc_gen.SendDispatchResponse
		expectedErr      error
	}{
		{
			name:         "failed: dispatch withsuch id does not exist",
			args:         arguments,
			mockedValues: &mocked{expectedSendErr: service_err.NotFoundErr},
			expectedErr:  status.Error(codes.NotFound, service_err.NotFoundErr.Error()),
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
			expectedResponse: &grpc_gen.SendDispatchResponse{},
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
		req *grpc_gen.GetAllDispatchesRequest
	}
	type mocked struct {
		expectedDispatches []deps.DispatchData
		expectedGetErr     error
	}

	dsMock := service_mock.NewDispatchService(t)
	loggerMock := logger_mock.NewLogger()
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
		req: &grpc_gen.GetAllDispatchesRequest{},
	}
	dispatches := []deps.DispatchData{{
		Id:                 "id",
		Label:              "label",
		CountOfSubscribers: 2,
	}}
	dispatchProtos := make([]*grpc_gen.DispatchData, 0, len(dispatches))
	for _, d := range dispatches {
		dispatchProtos = append(dispatchProtos, ToProtoDispatch(d))
	}

	tests := []struct {
		name             string
		args             *args
		mockedValues     *mocked
		expectedResponse *grpc_gen.GetAllDispatchesResponse
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
			expectedResponse: &grpc_gen.GetAllDispatchesResponse{
				Dispatches: []*grpc_gen.DispatchData{},
			},
		},
		{
			name:         "success: got all dispatches",
			args:         arguments,
			mockedValues: &mocked{expectedDispatches: dispatches},
			expectedResponse: &grpc_gen.GetAllDispatchesResponse{
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
				UnimplementedDispatchServiceServer: grpc_gen.UnimplementedDispatchServiceServer{},
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
