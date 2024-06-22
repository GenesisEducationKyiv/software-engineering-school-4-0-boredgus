package controllers

import (
	"context"
	"net/http"
	"strconv"
	client_mocks "subscription-api/internal/mocks/clients"
	controllers_mocks "subscription-api/internal/mocks/controllers"
	"subscription-api/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetExchangeRate_Controller(t *testing.T) {
	type mocked struct {
		ctx                    context.Context
		expectedRates          map[string]float64
		expectedConvertErr     error
		expectedResponseStatus int
		expectedResponseStr    string
	}

	csClientMock := client_mocks.NewCurrencyServiceClient(t)
	contextMock := controllers_mocks.NewContext(t)

	setup := func(m *mocked) func() {
		contextCall := contextMock.EXPECT().
			Context().Once().Return(m.ctx)
		convertCall := csClientMock.EXPECT().
			Convert(m.ctx, services.ConvertCurrencyParams{
				Base:   "USD",
				Target: []string{"UAH"},
			}).Once().NotBefore(contextCall).Return(m.expectedRates, m.expectedConvertErr)
		statusCall := contextMock.EXPECT().
			Status(m.expectedResponseStatus).NotBefore(convertCall).Maybe()
		stringCall := contextMock.EXPECT().
			String(m.expectedResponseStatus, m.expectedResponseStr).NotBefore(convertCall).Maybe()

		return func() {
			contextCall.Unset()
			convertCall.Unset()
			statusCall.Unset()
			stringCall.Unset()
		}
	}
	ctx := context.Background()
	uahRate := 30.0

	tests := []struct {
		name   string
		mocked *mocked
	}{
		{
			name: "failed: got error from currency service",
			mocked: &mocked{
				ctx:                    ctx,
				expectedConvertErr:     assert.AnError,
				expectedResponseStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "successfuly got exchange rate",
			mocked: &mocked{
				ctx:                    ctx,
				expectedRates:          map[string]float64{"UAH": uahRate},
				expectedResponseStr:    strconv.FormatFloat(uahRate, 'g', 7, 64),
				expectedResponseStatus: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked)
			defer cleanup()

			GetExchangeRate(contextMock, csClientMock)
		})
	}
}
