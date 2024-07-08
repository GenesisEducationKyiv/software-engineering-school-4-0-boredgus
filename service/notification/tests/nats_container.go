package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NATSContainer struct {
	Container        *nats.NATSContainer
	ConnectionString string
}

func CreateNatsContainer(ctx context.Context) (*NATSContainer, error) {
	natsContainer, err := nats.Run(ctx, "nats:2.10",
		nats.WithArgument("jetstream", "true"),
		nats.WithArgument("http_port", "8222"),
		nats.WithArgument("port", "4222"),
		testcontainers.WithWaitStrategyAndDeadline(
			5*time.Second,
			wait.ForLog("Listening for client connections on"),
		),
	)
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
