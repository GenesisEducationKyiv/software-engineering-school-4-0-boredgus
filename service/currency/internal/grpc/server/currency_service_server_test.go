package grpc_server

import (
	"context"
	"testing"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"
	logger_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/mocks/logger"
	service_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/mocks/service"
	service "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_CurrencyServiceServer_Convert(t *testing.T) {
	type args struct {
		req *grpc_gen.ConvertRequest
	}
	type mocked struct {
		rates      map[string]float64
		convertErr error
	}

	csMock := service_mock.NewCurrencyService(t)
	loggerMock := logger_mock.NewLogger()
	setup := func(m *mocked, args service.ConvertCurrencyParams) func() {
		csCall := csMock.EXPECT().Convert(mock.Anything, args).Return(m.rates, m.convertErr).Once()

		return func() {
			csCall.Unset()
		}
	}
	rates := map[string]float64{"UAH": 39.4347, "EUR": 0.9201}

	tests := []struct {
		name             string
		args             args
		mockedValues     mocked
		expectedResponse *grpc_gen.ConvertResponse
		expectedErr      error
	}{
		{
			name:             "failed: unsupported currency provided",
			args:             args{&grpc_gen.ConvertRequest{BaseCurrency: "123"}},
			mockedValues:     mocked{convertErr: service.InvalidArgumentErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.InvalidArgument, service.InvalidArgumentErr.Error()),
		},
		{
			name:             "failed: no target currency provided",
			args:             args{&grpc_gen.ConvertRequest{BaseCurrency: "USD"}},
			mockedValues:     mocked{convertErr: service.InvalidArgumentErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.InvalidArgument, service.InvalidArgumentErr.Error()),
		},
		{
			name:             "failed: failed precodition",
			args:             args{&grpc_gen.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH"}}},
			mockedValues:     mocked{convertErr: service.FailedPreconditionErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.FailedPrecondition, service.FailedPreconditionErr.Error()),
		},
		{
			name:             "failed: unknown error from Convert",
			args:             args{&grpc_gen.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH"}}},
			mockedValues:     mocked{convertErr: assert.AnError},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:             "successfully converted currency",
			args:             args{&grpc_gen.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH", "EUR"}}},
			mockedValues:     mocked{rates: rates},
			expectedResponse: &grpc_gen.ConvertResponse{BaseCurrency: "USD", Rates: rates},
			expectedErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(&tt.mockedValues, service.ConvertCurrencyParams{
				Base:   tt.args.req.BaseCurrency,
				Target: tt.args.req.TargetCurrencies,
			})
			defer cleanup()

			service := &currencyServiceServer{
				UnimplementedCurrencyServiceServer: grpc_gen.UnimplementedCurrencyServiceServer{},
				service:                            csMock,
				logger:                             loggerMock,
			}
			actualResponse, actualErr := service.Convert(context.Background(), tt.args.req)

			assert.Equal(t, tt.expectedResponse, actualResponse)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
