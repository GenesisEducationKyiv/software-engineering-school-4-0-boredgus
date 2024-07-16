package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func main() {
	env, err := config.Env()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess("transaction-manager"))

	// connection to NATS broker
	natsConnection, err := nats.Connect(
		env.BrokerURL,
		nats.Name("subscription-service"),
	)
	panicOnError(err, "failed to connect to broker")

	natsBroker, err := broker.NewNatsBroker(natsConnection, logger)
	panicOnError(err, "failed to create broker client")

	serverURL := fmt.Sprintf("%s:%s", env.TransactionManagerAddress, env.TransactionManagerPort)
	customerServiceURL := fmt.Sprintf("%s:%s", env.CustomerServiceAddress, env.CustomerServicePort)
	dispatchServiceURL := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)

	transactionManager := service.NewTransactionManager(
		serverURL,
		service.NewCustomerService(customerServiceURL),
		service.NewDispatchService(dispatchServiceURL),
		natsBroker,
	)
	server := grpc.NewServer()

	lis, err := net.Listen("tcp", serverURL)
	panicOnError(err, fmt.Sprintf("failed to listen %s", serverURL))

	logger.Infof("transaction manager started at %s", serverURL)
	panicOnError(server.Serve(lis), "failed to serve")
}
