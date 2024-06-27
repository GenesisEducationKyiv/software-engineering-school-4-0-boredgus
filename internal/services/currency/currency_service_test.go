package currency_service

import (
	"context"
	client_mocks "subscription-api/internal/mocks/clients"
	"subscription-api/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CurrencyService_Convert(t *testing.T) {
	type args struct {
		ctx    context.Context
		params services.ConvertCurrencyParams
	}
	type mocked struct {
		rates      map[string]float64
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
	rates := map[string]float64{"UAH": 30}

	tests := []struct {
		name          string
		args          args
		mockedValues  mocked
		expectedRates map[string]float64
		expectedErr   error
	}{
		{
			name:          "failed: no target currencies provided",
			args:          args{},
			expectedRates: nil,
			expectedErr:   services.InvalidArgumentErr,
		},
		{
			name: "failed: unsupported currency provided",
			args: args{
				params: services.ConvertCurrencyParams{Base: "invalid", Target: []string{"uah"}},
			},
			mockedValues:  mocked{convertErr: assert.AnError},
			expectedRates: nil,
			expectedErr:   services.InvalidArgumentErr,
		},
		{
			name: "successfuly converted currency",
			args: args{
				params: services.ConvertCurrencyParams{Base: "usd", Target: []string{"UAH"}},
			},
			mockedValues:  mocked{rates: rates},
			expectedRates: rates,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues)
			defer cleanup()

			service := &currencyService{
				currencyAPIClient: currencyAPIMock,
			}
			actualRates, actualErr := service.Convert(tt.args.ctx, tt.args.params)
			assert.Equal(t, tt.expectedRates, actualRates)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
