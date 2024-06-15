package cs

import (
	"context"
	"subscription-api/internal/entities"
	client_mocks "subscription-api/internal/mocks/clients"
	s "subscription-api/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CurrencyService_Convert(t *testing.T) {
	type args struct {
		ctx    context.Context
		params s.ConvertCurrencyParams
	}
	type mocked struct {
		rates      map[entities.Currency]float64
		convertErr error
	}
	currencyAPIMock := client_mocks.NewCurrencyAPIClient(t)
	setup := func(m mocked) func() {
		apiCall := currencyAPIMock.EXPECT().Convert(mock.Anything, mock.Anything, mock.Anything).
			Return(m.rates, m.convertErr)

		return func() {
			apiCall.Unset()
		}
	}
	rates := map[entities.Currency]float64{
		entities.UkrainianHryvnia: 30,
	}
	tests := []struct {
		name    string
		args    args
		mocked  mocked
		want    map[entities.Currency]float64
		wantErr error
	}{
		{
			name:    "no target currencies provided",
			args:    args{},
			want:    nil,
			wantErr: s.InvalidArgumentErr,
		},
		{
			name:    "success",
			args:    args{params: s.ConvertCurrencyParams{Target: []entities.Currency{entities.UkrainianHryvnia}}},
			mocked:  mocked{rates: rates},
			want:    rates,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clenup := setup(tt.mocked)
			defer clenup()

			s := &currencyService{
				currencyAPIClient: currencyAPIMock,
			}
			got, err := s.Convert(tt.args.ctx, tt.args.params)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
