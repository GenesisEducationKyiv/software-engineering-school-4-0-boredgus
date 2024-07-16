package clients

import (
	"context"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/clients/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if status.Code(err) == codes.AlreadyExists {
		return service.ErrAlreadyExists
	}

	return err
}
func (c *customerServiceClient) CreateCustomerRevert(ctx context.Context, email string) error {
	_, err := c.client.CreateCustomerRevert(ctx, &grpc_gen.CreateCustomerRequest{Email: email})

	return err
}
