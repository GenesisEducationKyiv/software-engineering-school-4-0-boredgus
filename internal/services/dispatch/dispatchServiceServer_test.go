package dispatch_service

import (
	"context"
	"fmt"
	config_mocks "subscription-api/internal/mocks/config"
	services_mocks "subscription-api/internal/mocks/services"
	"subscription-api/internal/services"
	grpc "subscription-api/internal/services/dispatch/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_DispatchServiceServer_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc.SubscribeForDispatchRequest
	}
	type mocked struct {
		subscribeErr error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger(t)
	internalError := fmt.Errorf("internal-error")
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			SubscribeForDispatch(a.ctx, a.req.Email, a.req.DispatchId).
			Once().Return(m.subscribeErr)
		logCall := loggerMock.EXPECT().
			Infof(mock.Anything, mock.Anything, mock.Anything)

		return func() {
			dsCall.Unset()
			logCall.Unset()
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
		name    string
		args    *args
		mocked  *mocked
		want    *grpc.SubscribeForDispatchResponse
		wantErr error
	}{
		{
			name:    "user already subscribed for this dispatch",
			args:    arguments,
			mocked:  &mocked{subscribeErr: services.UniqueViolationErr},
			wantErr: status.Error(codes.AlreadyExists, services.UniqueViolationErr.Error()),
		},
		{
			name:    "dispatch with such id does not exist",
			args:    arguments,
			mocked:  &mocked{subscribeErr: services.NotFoundErr},
			wantErr: status.Error(codes.NotFound, services.NotFoundErr.Error()),
		},
		{
			name:    "internal error occured",
			args:    arguments,
			mocked:  &mocked{subscribeErr: internalError},
			wantErr: status.Error(codes.Internal, internalError.Error()),
		},
		{
			name:   "success",
			args:   arguments,
			mocked: &mocked{},
			want:   &grpc.SubscribeForDispatchResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
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
		sendErr error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger(t)

	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			SendDispatch(a.ctx, a.req.DispatchId).
			Once().Return(m.sendErr)
		logCall := loggerMock.EXPECT().
			Infof(mock.Anything, mock.Anything, mock.Anything)

		return func() {
			dsCall.Unset()
			logCall.Unset()
		}
	}

	internalError := fmt.Errorf("internal-error")
	arguments := &args{
		ctx: context.Background(),
		req: &grpc.SendDispatchRequest{
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		want    *grpc.SendDispatchResponse
		wantErr error
	}{
		{
			name:    "dispatch not found",
			args:    arguments,
			mocked:  &mocked{sendErr: services.NotFoundErr},
			wantErr: status.Error(codes.NotFound, services.NotFoundErr.Error()),
		},
		{
			name:    "internal error occured",
			args:    arguments,
			mocked:  &mocked{sendErr: internalError},
			wantErr: status.Error(codes.Internal, internalError.Error()),
		},
		{
			name:   "success",
			args:   arguments,
			mocked: &mocked{},
			want:   &grpc.SendDispatchResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc.UnimplementedDispatchServiceServer{},
			}
			got, err := s.SendDispatch(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_DispatchServiceServer_GetAllDispatches(t *testing.T) {
	type args struct {
		ctx context.Context
		req *grpc.GetAllDispatchesRequest
	}
	type mocked struct {
		dispatches []services.DispatchData
		getErr     error
	}

	dsMock := services_mocks.NewDispatchService(t)
	loggerMock := config_mocks.NewLogger(t)
	internalError := fmt.Errorf("internal-error")
	setup := func(m *mocked, a *args) func() {
		dsCall := dsMock.EXPECT().
			GetAllDispatches(a.ctx).
			Once().Return(m.dispatches, m.getErr)
		logCall := loggerMock.EXPECT().
			Infof(mock.Anything, mock.Anything, mock.Anything)

		return func() {
			dsCall.Unset()
			logCall.Unset()
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
		dispatchProtos = append(dispatchProtos, d.ToProto())
	}

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		want    *grpc.GetAllDispatchesResponse
		wantErr error
	}{
		{
			name:    "internal err",
			args:    arguments,
			mocked:  &mocked{getErr: internalError},
			wantErr: status.Error(codes.Internal, internalError.Error()),
		},
		{
			name:   "there is no dispatches",
			args:   arguments,
			mocked: &mocked{},
			want: &grpc.GetAllDispatchesResponse{
				Dispatches: []*grpc.DispatchData{},
			},
		},
		{
			name:   "success",
			args:   arguments,
			mocked: &mocked{dispatches: dispatches},
			want: &grpc.GetAllDispatchesResponse{
				Dispatches: dispatchProtos,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				service:                            dsMock,
				logger:                             loggerMock,
				UnimplementedDispatchServiceServer: grpc.UnimplementedDispatchServiceServer{},
			}
			got, err := s.GetAllDispatches(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
