package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/clients/chain"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"
	grpc_server "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"

	"google.golang.org/grpc"
)

func must[T interface{}](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to gen envirinment variables")
	logger := config.InitLogger(env.Mode, config.WithProcess("currency-service"))
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

	// starting of server
	grpcServer := grpc.NewServer()
	grpc_gen.RegisterCurrencyServiceServer(grpcServer, currencyServiceServer)

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("currency service started at %s", url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}
