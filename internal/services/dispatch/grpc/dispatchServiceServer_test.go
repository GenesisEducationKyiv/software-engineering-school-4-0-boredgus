package grpc

import (
	"context"
	"fmt"
	e "subscription-api/internal/entities"
	config_mocks "subscription-api/internal/mocks/config"
	services_mocks "subscription-api/internal/mocks/services"
	"subscription-api/internal/services"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_DispatchServiceServer_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb_ds.SubscribeForDispatchRequest
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
		req: &pb_ds.SubscribeForDispatchRequest{
			Email:      "email",
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		want    *pb_ds.SubscribeForDispatchResponse
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
			want:   &pb_ds.SubscribeForDispatchResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				s:                                  dsMock,
				l:                                  loggerMock,
				UnimplementedDispatchServiceServer: pb_ds.UnimplementedDispatchServiceServer{},
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
		req *pb_ds.SendDispatchRequest
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
		req: &pb_ds.SendDispatchRequest{
			DispatchId: "dispatch-id",
		},
	}

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		want    *pb_ds.SendDispatchResponse
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
			want:   &pb_ds.SendDispatchResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				s:                                  dsMock,
				l:                                  loggerMock,
				UnimplementedDispatchServiceServer: pb_ds.UnimplementedDispatchServiceServer{},
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
		req *pb_ds.GetAllDispatchesRequest
	}
	type mocked struct {
		dispatches []e.CurrencyDispatch
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
		req: &pb_ds.GetAllDispatchesRequest{},
	}
	dispatches := []e.CurrencyDispatch{{
		Id:           "id",
		Label:        "label",
		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
		TemplateName: "template",
		Details: e.CurrencyDispatchDetails{
			BaseCurrency:     "base",
			TargetCurrencies: []string{"target"}},
		CountOfSubscribers: 2,
	}}
	dispatchProtos := make([]*pb_ds.DispatchData, 0, len(dispatches))
	for _, d := range dispatches {
		dispatchProtos = append(dispatchProtos, d.ToProto())
	}

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		want    *pb_ds.GetAllDispatchesResponse
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
			want: &pb_ds.GetAllDispatchesResponse{
				Dispatches: []*pb_ds.DispatchData{},
			},
		},
		{
			name:   "success",
			args:   arguments,
			mocked: &mocked{dispatches: dispatches},
			want: &pb_ds.GetAllDispatchesResponse{
				Dispatches: dispatchProtos,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchServiceServer{
				s:                                  dsMock,
				l:                                  loggerMock,
				UnimplementedDispatchServiceServer: pb_ds.UnimplementedDispatchServiceServer{},
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
