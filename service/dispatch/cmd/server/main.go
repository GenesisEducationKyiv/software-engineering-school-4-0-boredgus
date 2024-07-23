package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/server"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

const (
	MicroserviceName   string        = "dispatch-service"
	MetricPushInterval time.Duration = 15 * time.Second
)

func main() {
	env, err := config.Env()
	panicOnError(err, "failed to init environment variables")
	logger := config.InitLogger(env.Mode, config.WithProcess(MicroserviceName))
	defer logger.Flush()

	// connection to db
	postgresqlDB, err := db.NewPostrgreSQL(
		env.PostgreSQLConnString,
		db.PostgeSQLMigrationsUp,
	)
	panicOnError(err, "failed to connect to postgresql db")
	defer postgresqlDB.Close()

	storage := repo.NewStore(postgresqlDB, db.IsPqError)

	// initialization of dispatch service server
	dispatchService := service.NewDispatchService(
		repo.NewUserRepo(storage),
		repo.NewSubRepo(storage),
		repo.NewDispatchRepo(storage),
	)
	dispatchServiceServer := server.NewDispatchServiceServer(dispatchService, logger)

	// initialization of metrics interceptor
	commonMetricLabels := prometheus.Labels{"service": MicroserviceName}
	serverMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerCounterOptions(grpcprom.WithConstLabels(commonMetricLabels)),
		grpcprom.WithServerHandlingTimeHistogram(grpcprom.WithHistogramConstLabels(commonMetricLabels)),
	)

	// starting of grpc server
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(serverMetrics.UnaryServerInterceptor()),
	)
	grpc_gen.RegisterDispatchServiceServer(server, dispatchServiceServer)
	serverMetrics.InitializeMetrics(server)

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// schedulling of metrics push
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduleMetricsPush(ctx, env.MetricsGatewayURL, serverMetrics, logger)

	logger.Infof("dispatch service started at %s", url)
	panicOnError(server.Serve(lis), "failed to serve")
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
