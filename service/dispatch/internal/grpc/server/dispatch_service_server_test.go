package server

import (
	"context"
	"testing"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"

	logger_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/logger"
	service_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		want         *emptypb.Empty
		wantErr      error
	}{
		{
			name:         "failed: user already subscribed for this dispatch",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: service.UniqueViolationErr},
			wantErr:      status.Error(codes.AlreadyExists, service.UniqueViolationErr.Error()),
		},
		{
			name:         "failed: dispatch with such id does not exist",
			args:         arguments,
			mockedValues: &mocked{expectedSubscribeErr: service.NotFoundErr},
			wantErr:      status.Error(codes.NotFound, service.NotFoundErr.Error()),
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
			want:         &emptypb.Empty{},
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
