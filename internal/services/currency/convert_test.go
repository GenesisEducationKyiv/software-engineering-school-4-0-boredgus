package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"subscription-api/internal/entities"
	"subscription-api/pkg/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CurrencyService_Convert(t *testing.T) {
	invalidCurrency := entities.Currency("invalid-currency")
	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/latest/%s", invalidCurrency) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"result":"error","error-type":"%s"}`, InvalidArgumentErr)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","conversion_rates":{"USD":1,"EUR":0.9201,"GBP":0.7883,"PLN":3.9255,"UAH":39.4347}}`))
	}))
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest/USD" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"result":"error","error-type":"some-error"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`}`))
	}))
	defer validServer.Close()
	defer invalidServer.Close()

	type fields struct {
		APIBasePath string
	}
	type args struct {
		ctx    context.Context
		params ConvertParams
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[entities.Currency]float64
		wantErr error
	}{
		{
			name:    "failed to fetch exchange rate info from thrird-party api",
			fields:  fields{},
			args:    args{ctx: ctx, params: ConvertParams{}},
			wantErr: InvalidArgumentErr,
		},
		{
			name:    "failed to fetch exchange rate info from thrird-party api",
			fields:  fields{APIBasePath: "invalid-url"},
			args:    args{ctx: ctx, params: ConvertParams{To: []entities.Currency{"UAH"}}},
			wantErr: InvalidRequestErr,
		},
		{
			name:    "invalid format of thrird-party api response",
			fields:  fields{APIBasePath: invalidServer.URL},
			args:    args{ctx: ctx, params: ConvertParams{To: []entities.Currency{"UAH"}}},
			wantErr: tools.ParseErr,
		},
		{
			name:    "unsupported currency provided",
			fields:  fields{APIBasePath: validServer.URL},
			args:    args{ctx: ctx, params: ConvertParams{From: invalidCurrency, To: []entities.Currency{"UAH"}}},
			wantErr: InvalidArgumentErr,
		},
		{
			name:    "unexpected thrird-party api response",
			fields:  fields{APIBasePath: invalidServer.URL},
			args:    args{ctx: ctx, params: ConvertParams{From: "USD", To: []entities.Currency{"UAH"}}},
			wantErr: FailedPreconditionErr,
		},
		{
			name:   "succesfully converted",
			fields: fields{APIBasePath: validServer.URL},
			args:   args{ctx: ctx, params: ConvertParams{From: "USD", To: []entities.Currency{"UAH", "EUR"}}},
			want:   map[entities.Currency]float64{"EUR": 0.9201, "UAH": 39.4347},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := currencyService{
				APIBasePath: tt.fields.APIBasePath,
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
