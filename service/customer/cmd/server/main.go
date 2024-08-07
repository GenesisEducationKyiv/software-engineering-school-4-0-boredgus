package main

import (
	"context"
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/db"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess(env.MicroserviceName))

	database, err := db.NewDatabase(env.DatabaseURL, env.DatabaseSchema)
	panicOnError(err, "failed to setup database")

	customerService := service.NewCustomerService(
		repo.NewCustomerRepo(database),
	)
	customerServer := server.NewCustomerServiceServer(customerService, logger)

	// initialization of metrics interceptor
	commonMetricLabels := prometheus.Labels{"service": env.MicroserviceName}
	serverMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerCounterOptions(grpcprom.WithConstLabels(commonMetricLabels)),
		grpcprom.WithServerHandlingTimeHistogram(grpcprom.WithHistogramConstLabels(commonMetricLabels)),
	)

	// initialization of grpc server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(serverMetrics.UnaryServerInterceptor()),
	)
	grpc_gen.RegisterCustomerServiceServer(grpcServer, customerServer)
	serverMetrics.InitializeMetrics(grpcServer)

	// expose metrics
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(serverMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go metrics.
		NewMetricsServer(logger, ":"+env.MetricsPort, env.MetricsRoute, promRegistry).
		Run(ctx)
	url := fmt.Sprintf("%s:%s", env.CustomerServiceAddress, env.CustomerServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// schedulling of metrics push
	go metrics.NewMetricsPusher(logger).Push(ctx, metrics.PushParams{
		URLToFetchMetrics: fmt.Sprintf("http://localhost:%v%v", env.MetricsPort, env.MetricsRoute),
		URLToPushMetrics:  env.MetricsGatewayURL,
		PushInterval:      metrics.DefaultMetricsPushInterval,
	})

	logger.Infof("%s started at %s", env.MicroserviceName, url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}
