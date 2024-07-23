package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients/chain"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"
	grpc_server "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
)

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

const (
	MicroserviceName   string        = "currency-service"
	MetricPushInterval time.Duration = 15 * time.Second
)

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to gen envirinment variables")
	logger := config.InitLogger(env.Mode, config.WithProcess(MicroserviceName))
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
	commonMetricLabels := prometheus.Labels{"service": MicroserviceName}
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

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	// schedulling of metrics push
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduleMetricsPush(ctx, env.MetricsGatewayURL, serverMetrics, logger)

	// start of the server
	logger.Infof("currency service started at %s", url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}

func scheduleMetricsPush(ctx context.Context, urlToPush string, collector prometheus.Collector, logger config.Logger) {
	pusher := push.New(urlToPush, MicroserviceName).Collector(collector)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			time.Sleep(MetricPushInterval)

			if err := pusher.Push(); err != nil {
				logger.Errorf("failed to push metrics: %v", err)
			}
		}
	}
}
