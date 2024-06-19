package controllers

import (
	"context"
	"fmt"
	"net/http"
	"subscription-api/internal/entities"
	client_mocks "subscription-api/internal/mocks/clients"
	controllers_mocks "subscription-api/internal/mocks/controllers"
	"subscription-api/internal/services"
	"testing"
)

func Test_GetExchangeRate_Controller(t *testing.T) {
	type mocked struct {
		ctx            context.Context
		rates          map[string]float64
		convertErr     error
		responseStatus int
		responseStr    string
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
			}).Once().NotBefore(contextCall).Return(m.rates, m.convertErr)
		statusCall := contextMock.EXPECT().
			Status(m.responseStatus).NotBefore(convertCall).Maybe()
		stringCall := contextMock.EXPECT().
			String(m.responseStatus, m.responseStr).NotBefore(convertCall).Maybe()

		return func() {
			contextCall.Unset()
			convertCall.Unset()
			statusCall.Unset()
			stringCall.Unset()
		}
	}
	ctx := context.Background()
	uahRate := 30

	tests := []struct {
		name   string
		mocked *mocked
	}{
		{
			name: "failed to convert currency",
			mocked: &mocked{
				ctx:            ctx,
				convertErr:     fmt.Errorf("some-err"),
				responseStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "successfuly got exchange rate",
			mocked: &mocked{
				ctx: ctx,
				rates: map[string]float64{
					string(entities.UkrainianHryvnia): float64(uahRate)},
				responseStr:    strconv.Itoa(uahRate),
				responseStatus: http.StatusOK,
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
