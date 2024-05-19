package grpc

import (
	"context"
	"fmt"
	"subscription-api/internal/entities"
	mocks "subscription-api/internal/mocks/services"
	cs "subscription-api/internal/services/currency"
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
	csMock := mocks.NewCurrencyService(t)
	internalError := fmt.Errorf("internal-error")
	setup := func(res *mockedRes, args cs.ConvertParams) func() {
		csCall := csMock.EXPECT().Convert(mock.Anything, args).Return(res.convertedRates, res.convertErr).Once()
		return func() {
			csCall.Unset()
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
			mockedRes: mockedRes{convertErr: cs.InvalidArgumentErr},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, cs.InvalidArgumentErr.Error()),
		},
		{
			name:      "no target currencies provided",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "123"}},
			mockedRes: mockedRes{convertErr: cs.InvalidArgumentErr},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, cs.InvalidArgumentErr.Error()),
		},
		{
			name:      "failed precodition",
			args:      args{&pb_cs.ConvertRequest{BaseCurrency: "USD"}},
			mockedRes: mockedRes{convertErr: cs.FailedPreconditionErr},
			want:      nil,
			wantErr:   status.Error(codes.FailedPrecondition, cs.FailedPreconditionErr.Error()),
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
			reset := setup(&tt.mockedRes, cs.ConvertParams{
				From: entities.Currency(tt.args.req.BaseCurrency),
				To:   entities.FromString(tt.args.req.TargetCurrencies)})
			defer reset()

			s := &currencyServiceServer{
				UnimplementedCurrencyServiceServer: pb_cs.UnimplementedCurrencyServiceServer{},
				s:                                  csMock,
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
