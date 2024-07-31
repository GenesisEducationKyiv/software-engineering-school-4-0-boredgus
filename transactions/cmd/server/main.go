package main

import (
	"context"
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/clients"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env, err := config.Env()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess(env.MicroserviceName))

	// connection to NATS broker
	natsConnection, err := nats.Connect(
		env.BrokerURL,
		nats.Name(env.MicroserviceName),
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

	// initialization of metrics interceptor
	commonMetricLabels := prometheus.Labels{"service": env.MicroserviceName}
	serverMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerCounterOptions(grpcprom.WithConstLabels(commonMetricLabels)),
		grpcprom.WithServerHandlingTimeHistogram(grpcprom.WithHistogramConstLabels(commonMetricLabels)),
	)

	// initialization of grpc server
	transactionManagerServer := server.NewTransactionManagerServer(transactionManager, logger)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(serverMetrics.UnaryServerInterceptor()),
	)
	grpc_gen.RegisterTransactionManagerServer(server, transactionManagerServer)
	serverMetrics.InitializeMetrics(server)

	// expose metrics
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(serverMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go metrics.
		NewMetricsServer(logger, ":"+env.MetricsPort, env.MetricsRoute, promRegistry).
		Run(ctx)

	lis, err := net.Listen("tcp", serverURL)
	panicOnError(err, fmt.Sprintf("failed to listen %s", serverURL))

	// schedulling of metrics push
	go metrics.NewMetricsPusher(logger).
		Push(ctx, metrics.PushParams{
			URLToFetchMetrics: fmt.Sprintf("http://localhost:%v%v", env.MetricsPort, env.MetricsRoute),
			URLToPushMetrics:  env.MetricsGatewayURL,
			PushInterval:      metrics.DefaultMetricsPushInterval,
		})

	// start of the server
	logger.Infof("%s started at %s", env.MicroserviceName, serverURL)
	panicOnError(server.Serve(lis), "failed to serve")
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}
