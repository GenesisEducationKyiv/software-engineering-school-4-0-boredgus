package grpc

import (
	"context"
	"fmt"
	"subscription-api/internal/entities"
	config_mocks "subscription-api/internal/mocks/config"
	s_mocks "subscription-api/internal/mocks/services"
	s "subscription-api/internal/services"
	pb_cs "subscription-api/pkg/grpc/currency_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_CurrencyServiceServer_Convert(t *testing.T) {
	type args struct {
		req *pb_cs.ConvertRequest
	}
	type mockedRes struct {
		convertedRates map[entities.Currency]float64
		convertErr     error
	}
	csMock := s_mocks.NewCurrencyService(t)
	loggerMock := config_mocks.NewLogger(t)
	internalError := fmt.Errorf("internal-error")
	setup := func(res *mockedRes, args s.ConvertCurrencyParams) func() {
		csCall := csMock.EXPECT().Convert(mock.Anything, args).Return(res.convertedRates, res.convertErr).Once()
		logCall := loggerMock.EXPECT().Infof(mock.Anything, mock.Anything, mock.Anything)

		return func() {
			csCall.Unset()
			logCall.Unset()
		}
	}
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		want      *pb_cs.ConvertResponse
		wantErr   error
	}{
		{
			name:      "unsupported currency provided",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "123"}},
			mockedRes: mockedRes{convertErr: s.InvalidArgumentErr},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, s.InvalidArgumentErr.Error()),
		},
		{
			name:      "no target currencies provided",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "123"}},
			mockedRes: mockedRes{convertErr: s.InvalidArgumentErr},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, s.InvalidArgumentErr.Error()),
		},
		{
			name:      "failed precodition",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "USD"}},
			mockedRes: mockedRes{convertErr: s.FailedPreconditionErr},
			want:      nil,
			wantErr:   status.Error(codes.FailedPrecondition, s.FailedPreconditionErr.Error()),
		},
		{
			name:      "internal error",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "USD"}},
			mockedRes: mockedRes{convertErr: internalError},
			want:      nil,
			wantErr:   status.Error(codes.Internal, internalError.Error()),
		},
		{
			name: "successfully converted",
			args: args{&pb_cs.ConvertRequest{BaseCurrency: "USD", TargetCurrencies: []string{"UAH", "EUR"}}},
			mockedRes: mockedRes{
				convertedRates: map[entities.Currency]float64{"UAH": 39.4347, "EUR": 0.9201}},
			want: &pb_cs.ConvertResponse{
				BaseCurrency: "USD",
				Rates:        map[string]float64{"UAH": 39.4347, "EUR": 0.9201}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset := setup(&tt.mockedRes, s.ConvertCurrencyParams{
				Base:   entities.Currency(tt.args.req.BaseCurrency),
				Target: entities.FromString(tt.args.req.TargetCurrencies)})
			defer reset()

			s := &currencyServiceServer{
				UnimplementedCurrencyServiceServer: pb_cs.UnimplementedCurrencyServiceServer{},
				s:                                  csMock,
				l:                                  loggerMock,
			}
			got, err := s.Convert(context.Background(), tt.args.req)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
