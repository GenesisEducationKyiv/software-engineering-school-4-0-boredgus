package main

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/currency"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/metrics"
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
	panicOnError(err, "failed to get environment variables")

	logger := logger.InitLogger(env.Mode, logger.WithProcess(env.MicroserviceName))
	defer logger.Flush()

	// initialization of service clients
	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	transactionManagerConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.TransactionManagerAddress, env.TransactionManagerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to dispatch service grpc server")
	defer transactionManagerConn.Close()

	// schedulling of metrics push
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go metrics.NewMetricsPusher(logger).
		Push(ctx, metrics.PushParams{
			URLToFetchMetrics: fmt.Sprintf("http://localhost:%v%v", env.Port, config.MetricsURL),
			URLToPushMetrics:  env.MetricsGatewayURL,
			PushInterval:      metrics.DefaultMetricsPushInterval,
		})

	// initialization of router
	router := config.GetRouter(&config.APIParams{
		CurrencyService:  currency.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService:  dispatch.NewTransactionManagerClient(transactionManagerConn),
		Logger:           logger,
		MicroserviceName: env.MicroserviceName,
	})

	logger.Infof("started %s at %v port", env.MicroserviceName, env.Port)

	panicOnError(router.Run(":"+env.Port), "failed to start server")
}
