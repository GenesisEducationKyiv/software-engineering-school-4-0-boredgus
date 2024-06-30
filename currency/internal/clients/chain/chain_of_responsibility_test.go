package chain

import (
	"context"
	"testing"

	client_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/mocks/client"
	"github.com/stretchr/testify/assert"
)

func Test_CurrencyAPIChain_Convert(t *testing.T) {
	type args struct {
		ctx         context.Context
		baseCcy     string
		targetCcies []string
	}
	type mocked struct {
		primaryAPIRates   map[string]float64
		primaryAPIErr     error
		secondaryAPIRates map[string]float64
		secondaryAPIErr   error
	}

	primaryAPIMock := client_mock.NewCurrencyAPIClient(t)
	secondaryAPIMock := client_mock.NewCurrencyAPIClient(t)
	secondaryChain := NewCurrencyAPIChain(secondaryAPIMock)

	setup := func(m *mocked, a *args) func() {
		primaryCall := primaryAPIMock.EXPECT().
			Convert(a.ctx, a.baseCcy, a.targetCcies).Once().
			Return(m.primaryAPIRates, m.primaryAPIErr)
		secondaryCall := secondaryAPIMock.EXPECT().
			Convert(a.ctx, a.baseCcy, a.targetCcies).Maybe().NotBefore(primaryCall).
			Return(m.secondaryAPIRates, m.secondaryAPIErr)

		return func() {
			primaryCall.Unset()
			secondaryCall.Unset()
		}
	}

	arguments := &args{
		ctx:         context.Background(),
		baseCcy:     "usd",
		targetCcies: []string{"uah", "eur"},
	}

	rates := map[string]float64{"uah": 40.0, "eur": 0.95}

	tests := []struct {
		name          string
		mockedValues  *mocked
		expectedRates map[string]float64
		expectedErr   error
	}{
		{
			name: "success: primary API converted currency",
			mockedValues: &mocked{
				primaryAPIRates: rates,
			},
			expectedRates: rates,
		},
		{
			name: "success: secondary API converted currency",
			mockedValues: &mocked{
				primaryAPIErr:     assert.AnError,
				secondaryAPIRates: rates,
			},
			expectedRates: rates,
		},
		{
			name: "failed: all APIs have failed to convert currency",
			mockedValues: &mocked{
				primaryAPIErr:   assert.AnError,
				secondaryAPIErr: assert.AnError,
			},
			expectedErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, arguments)
			defer cleanup()

			ch := &currencyAPIChain{api: primaryAPIMock}
			ch.SetNext(secondaryChain)

			actualRates, actualErr := ch.Convert(arguments.ctx, arguments.baseCcy, arguments.targetCcies)

			assert.Equal(t, tt.expectedRates, actualRates)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.NoError(t, actualErr)
		})
	}
}
