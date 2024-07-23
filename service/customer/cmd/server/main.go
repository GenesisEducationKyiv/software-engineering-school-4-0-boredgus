package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/db"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"google.golang.org/grpc"
)

const (
	MicroserviceName   string        = "customer-service"
	MetricPushInterval time.Duration = 15 * time.Second
)

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess(MicroserviceName))

	database, err := db.NewDatabase(env.DatabaseURL, env.DatabaseSchema)
	panicOnError(err, "failed to setup database")

	customerService := service.NewCustomerService(
		repo.NewCustomerRepo(database),
	)
	customerServer := server.NewCustomerServiceServer(customerService, logger)

	// initialization of metrics interceptor
	commonMetricLabels := prometheus.Labels{"service": MicroserviceName}
	serverMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerCounterOptions(grpcprom.WithConstLabels(commonMetricLabels)),
		grpcprom.WithServerHandlingTimeHistogram(grpcprom.WithHistogramConstLabels(commonMetricLabels)),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(serverMetrics.UnaryServerInterceptor()),
	)
	grpc_gen.RegisterCustomerServiceServer(grpcServer, customerServer)
	serverMetrics.InitializeMetrics(grpcServer)

	url := fmt.Sprintf("%s:%s", env.CustomerServiceAddress, env.CustomerServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// schedulling of metrics push
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduleMetricsPush(ctx, env.MetricsGatewayURL, serverMetrics, logger)

	logger.Infof("customer service started at %s", url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func scheduleMetricsPush(ctx context.Context, urlToPush string, collector prometheus.Collector, logger config.Logger) {
	pusher := push.New(urlToPush, MicroserviceName).Collector(collector)
	ticker := time.NewTicker(MetricPushInterval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			if err := pusher.Push(); err != nil {
				logger.Errorf("failed to push metrics: %v", err)
			}
		}
	}
}
