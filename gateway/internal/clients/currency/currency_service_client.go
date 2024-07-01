package currency

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/currency/gen"
	"google.golang.org/grpc"
)

type (
	ConvertCurrencyParams struct {
		Base   string
		Target []string
	}

	currencyServiceClient struct {
		client grpc_gen.CurrencyServiceClient
	}
)

func NewCurrencyServiceClient(conn grpc.ClientConnInterface) *currencyServiceClient {
	return &currencyServiceClient{
		client: grpc_gen.NewCurrencyServiceClient(conn),
	}
}

func (c *currencyServiceClient) Convert(ctx context.Context, base string, target []string) (map[string]float64, error) {
	resp, err := c.client.Convert(ctx, &grpc_gen.ConvertRequest{
		BaseCurrency:     base,
		TargetCurrencies: target,
	})
	if err != nil {
		return nil, err
	}

	return resp.Rates, nil
}
