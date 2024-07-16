package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/clients"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	customerServiceURL := fmt.Sprintf("%s:%s", env.CustomerServiceAddress, env.CustomerServicePort)
	dispatchServiceURL := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)

	// connection to gRPC server of customer service
	customerServiceConn, err := grpc.NewClient(
		customerServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to customer service grpc server")
	defer customerServiceConn.Close()

	// connection to gRPC server of customer service
	dispatchServiceConn, err := grpc.NewClient(
		dispatchServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to customer service grpc server")
	defer dispatchServiceConn.Close()

	// initialization of transaction manager
	serverURL := fmt.Sprintf("%s:%s", env.TransactionManagerAddress, env.TransactionManagerPort)
	transactionManager := service.NewTransactionManager(
		serverURL,
		clients.NewCustomerServiceClient(customerServiceConn),
		clients.NewDispatchServiceClient(dispatchServiceConn),
		broker.NewEventBroker(natsBroker, logger),
		logger,
	)
	transactionManagerServer := server.NewTransactionManagerServer(transactionManager, logger)
	server := grpc.NewServer()
	grpc_gen.RegisterTransactionManagerServer(server, transactionManagerServer)

	lis, err := net.Listen("tcp", serverURL)
	panicOnError(err, fmt.Sprintf("failed to listen %s", serverURL))

	logger.Infof("transaction manager started at %s", serverURL)
	panicOnError(server.Serve(lis), "failed to serve")
}
