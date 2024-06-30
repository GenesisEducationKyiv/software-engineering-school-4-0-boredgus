package client

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"
	"google.golang.org/grpc"
)

type currencyServiceClient struct {
	client grpc_gen.CurrencyServiceClient
}

func NewCurrencyServiceClient(conn grpc.ClientConnInterface) *currencyServiceClient {
	return &currencyServiceClient{
		client: grpc_gen.NewCurrencyServiceClient(conn),
	}
}

func (c *currencyServiceClient) Convert(ctx context.Context, params service.ConvertCurrencyParams) (map[string]float64, error) {
	resp, err := c.client.Convert(ctx, &grpc_gen.ConvertRequest{
		BaseCurrency:     params.Base,
		TargetCurrencies: params.Target,
	})
	if err != nil {
		return nil, err
	}

	return resp.Rates, nil
}
