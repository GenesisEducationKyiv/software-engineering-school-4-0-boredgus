package main

import (
	"context"
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/prometheus/client_golang/prometheus"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/server"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

const (
	microserviceName string = "dispatch-service"
	metricsPort      string = "8012"
	metricsPath      string = "/metrics"
)

func main() {
	env, err := config.Env()
	panicOnError(err, "failed to init environment variables")
	logger := config.InitLogger(env.Mode, config.WithProcess(microserviceName))
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
	commonMetricLabels := prometheus.Labels{"service": microserviceName}
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

	// expose metrics
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(serverMetrics)
	go metrics.ExposeMetrics(":"+metricsPort, metricsPath, promRegistry)

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// scheduling of metrics push
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go metrics.NewMetricsPusher(logger).
		Push(ctx, metrics.PushParams{
			URLToFetchMetrics: fmt.Sprintf("http://localhost:%v%v", metricsPort, metricsPath),
			URLToPushMetrics:  env.MetricsGatewayURL,
			PushInterval:      metrics.DefaultMetricsPushInterval,
		})

	logger.Infof("%s started at %s", microserviceName, url)
	panicOnError(server.Serve(lis), "failed to serve")
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}
