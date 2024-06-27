package currency_server

import (
	"context"
	config_mocks "subscription-api/internal/mocks/config"
	s_mocks "subscription-api/internal/mocks/services"
	s "subscription-api/internal/services"
	grpc "subscription-api/internal/services/currency/server/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_CurrencyServiceServer_Convert(t *testing.T) {
	type args struct {
		req *grpc.ConvertRequest
	}
	type mocked struct {
		rates      map[string]float64
		convertErr error
	}

	csMock := s_mocks.NewCurrencyService(t)
	loggerMock := config_mocks.NewLogger()
	setup := func(m *mocked, args s.ConvertCurrencyParams) func() {
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
		expectedResponse *grpc.ConvertResponse
		expectedErr      error
	}{
		{
			name:             "failed: unsupported currency provided",
			args:             args{&grpc.ConvertRequest{BaseCurrency: "123"}},
			mockedValues:     mocked{convertErr: s.InvalidArgumentErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.InvalidArgument, s.InvalidArgumentErr.Error()),
		},
		{
			name:             "failed: no target currency provided",
			args:             args{&grpc.ConvertRequest{BaseCurrency: "USD"}},
			mockedValues:     mocked{convertErr: s.InvalidArgumentErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.InvalidArgument, s.InvalidArgumentErr.Error()),
		},
		{
			name:             "failed: failed precodition",
			args:             args{&grpc.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH"}}},
			mockedValues:     mocked{convertErr: s.FailedPreconditionErr},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.FailedPrecondition, s.FailedPreconditionErr.Error()),
		},
		{
			name:             "failed: unknown error from Convert",
			args:             args{&grpc.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH"}}},
			mockedValues:     mocked{convertErr: assert.AnError},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Internal, assert.AnError.Error()),
		},
		{
			name:             "successfully converted currency",
			args:             args{&grpc.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH", "EUR"}}},
			mockedValues:     mocked{rates: rates},
			expectedResponse: &grpc.ConvertResponse{BaseCurrency: "USD", Rates: rates},
			expectedErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(&tt.mockedValues, s.ConvertCurrencyParams{
				Base:   tt.args.req.BaseCurrency,
				Target: tt.args.req.TargetCurrencies,
			})
			defer cleanup()

			service := &currencyServiceServer{
				UnimplementedCurrencyServiceServer: grpc.UnimplementedCurrencyServiceServer{},
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
