package grpc

import (
	"context"
	"subscription-api/internal/services"
	grpc_client "subscription-api/internal/services/currency/grpc"

	"google.golang.org/grpc"
)

type CurrencyServiceClient interface {
	Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error)
}

type currencyServiceClient struct {
	client grpc_client.CurrencyServiceClient
}

func NewCurrencyServiceClient(conn grpc.ClientConnInterface) *currencyServiceClient {
	return &currencyServiceClient{
		client: grpc_client.NewCurrencyServiceClient(conn),
	}
}

func (c *currencyServiceClient) Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error) {
	resp, err := c.client.Convert(ctx, &grpc_client.ConvertRequest{
		BaseCurrency:     params.Base,
		TargetCurrencies: params.Target,
	})
	if err != nil {
		return nil, err
	}

	return resp.Rates, nil
}
