package cs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	e "subscription-api/internal/entities"
	"subscription-api/internal/services"
	"subscription-api/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CurrencyService_Convert(t *testing.T) {
	invalidCurrency := e.Currency("invalid-currency")
	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/latest/%s", invalidCurrency) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(fmt.Sprintf(`{"result":"error","error-type":"%s"}`, services.InvalidArgumentErr)))
			require.NoError(t, err)

			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"result":"success","conversion_rates":{"USD":1,"EUR":0.9201,"GBP":0.7883,"PLN":3.9255,"UAH":39.4347}}`))
		require.NoError(t, err)
	}))
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest/USD" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"result":"error","error-type":"some-error"}`))
			require.NoError(t, err)

			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`}`))
		require.NoError(t, err)
	}))
	defer validServer.Close()
	defer invalidServer.Close()

	type fields struct {
		APIBasePath string
	}
	type args struct {
		ctx    context.Context
		params ConvertCurrencyParams
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[e.Currency]float64
		wantErr error
	}{
		{
			name:    "failed to fetch exchange rate info from thrird-party api",
			fields:  fields{},
			args:    args{ctx: ctx, params: ConvertCurrencyParams{}},
			wantErr: services.InvalidArgumentErr,
		},
		// {
		// 	name:    "failed to fetch exchange rate info from thrird-party api",
		// 	fields:  fields{APIBasePath: "invalid-url"},
		// 	args:    args{ctx: ctx, params: ConvertParams{To: []e.Currency{"UAH"}}},
		// 	wantErr: services.InvalidRequestErr,
		// },
		{
			name:    "invalid format of thrird-party api response",
			fields:  fields{APIBasePath: invalidServer.URL},
			args:    args{ctx: ctx, params: ConvertCurrencyParams{Target: []e.Currency{"UAH"}}},
			wantErr: utils.ParseErr,
		},
		{
			name:    "unsupported currency provided",
			fields:  fields{APIBasePath: validServer.URL},
			args:    args{ctx: ctx, params: ConvertCurrencyParams{Base: invalidCurrency, Target: []e.Currency{"UAH"}}},
			wantErr: services.InvalidArgumentErr,
		},
		{
			name:    "unexpected thrird-party api response",
			fields:  fields{APIBasePath: invalidServer.URL},
			args:    args{ctx: ctx, params: ConvertCurrencyParams{Base: "USD", Target: []e.Currency{"UAH"}}},
			wantErr: services.FailedPreconditionErr,
		},
		{
			name:   "succesfully converted",
			fields: fields{APIBasePath: validServer.URL},
			args:   args{ctx: ctx, params: ConvertCurrencyParams{Base: "USD", Target: []e.Currency{"UAH", "EUR"}}},
			want:   map[e.Currency]float64{"EUR": 0.9201, "UAH": 39.4347},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := currencyService{
				// currencyAPIClient: ,
			}
			got, err := e.Convert(tt.args.ctx, tt.args.params)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
