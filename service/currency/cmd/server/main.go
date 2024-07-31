package main

import (
	"context"
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients/chain"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"
	grpc_server "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to gen envirinment variables")
	logger := config.InitLogger(env.Mode, config.WithProcess(env.MicroserviceName))
	defer logger.Flush()

	// initialization of currency API clients
	httpClient := clients.NewHTTPClient()
	exchangeRateAPIClient := chain.NewCurrencyAPIChain(
		clients.NewExchangeRateAPIClient(httpClient, env.ExchangeCurrencyAPIKey, logger),
	)
	currencyBeaconAPIClient := chain.NewCurrencyAPIChain(
		clients.NewCurrencyBeaconAPIClient(httpClient, env.CurrencyBeaconAPIKey, logger),
	)
	freeCurrencyAPIClient := chain.NewCurrencyAPIChain(
		clients.NewFreeCurrencyAPIClient(httpClient, logger),
	)
	exchangeRateAPIClient.SetNext(currencyBeaconAPIClient)
	currencyBeaconAPIClient.SetNext(freeCurrencyAPIClient)

	// initalization currency service server
	currencyService := service.NewCurrencyService(exchangeRateAPIClient)
	currencyServiceServer := grpc_server.NewCurrencyServiceServer(
		currencyService,
		logger,
	)

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
	grpc_gen.RegisterCurrencyServiceServer(grpcServer, currencyServiceServer)
	serverMetrics.InitializeMetrics(grpcServer)

	// expose metrics
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(serverMetrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go metrics.
		NewMetricsServer(logger, ":"+env.MetricsPort, env.MetricsRoute, promRegistry).
		Run(ctx)

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// scheduling of metrics push
	go metrics.NewMetricsPusher(logger).
		Push(ctx, metrics.PushParams{
			URLToFetchMetrics: fmt.Sprintf("http://localhost:%v%v", env.MetricsPort, env.MetricsRoute),
			URLToPushMetrics:  env.MetricsGatewayURL,
			PushInterval:      metrics.DefaultMetricsPushInterval,
		})

	// start of the server
	logger.Infof("%s started at %s", env.MicroserviceName, url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}
