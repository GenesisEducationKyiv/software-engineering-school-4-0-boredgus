package clients

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/clients/gen"
	"google.golang.org/grpc"
)

type customerServiceClient struct {
	client grpc_gen.CustomerServiceClient
}

func NewCustomerServiceClient(conn grpc.ClientConnInterface) *customerServiceClient {
	return &customerServiceClient{
		client: grpc_gen.NewCustomerServiceClient(conn),
	}
}

func (c *customerServiceClient) CreateCustomer(ctx context.Context, email string) error {
	_, err := c.client.CreateCustomer(ctx, &grpc_gen.CreateCustomerRequest{Email: email})

	return err
}
func (c *customerServiceClient) CreateCustomerRevert(ctx context.Context, email string) error {
	_, err := c.client.CreateCustomerRevert(ctx, &grpc_gen.CreateCustomerRequest{Email: email})

	return err
}
