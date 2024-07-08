package tests

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/nats"
)

type NATSContainer struct {
	Container        *nats.NATSContainer
	ConnectionString string
}

func CreateNatsContainer(ctx context.Context) (*NATSContainer, error) {
	natsContainer, err := nats.Run(ctx, "nats:2.10")
	if err != nil {
		return nil, fmt.Errorf("failed to run NATS container: %w", err)
	}

	connString, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get NATS connection string: %w", err)
	}

	return &NATSContainer{
		Container:        natsContainer,
		ConnectionString: connString,
	}, nil
}
