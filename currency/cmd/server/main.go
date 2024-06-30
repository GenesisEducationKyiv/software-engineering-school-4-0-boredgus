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

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/utils"
	"google.golang.org/grpc"
)

func main() {
	env := utils.Must(config.GetEnv())
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
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("currency service started at %s", url)
	utils.PanicOnError(grpcServer.Serve(lis), "failed to serve")
}
