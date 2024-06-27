package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/currency-service/internal"
	"subscription-api/config"
	"subscription-api/internal/clients"
	currency_client "subscription-api/internal/clients/currency"
	currency_service "subscription-api/internal/services/currency"
	currency_grpc "subscription-api/internal/services/currency/grpc"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(config.ProdMode, config.WithProcess("currency-service"))
	defer logger.Flush()

	// initialization of currency API clients
	httpClient := clients.NewHTTPClient()
	exchangeRateAPIClient := currency_client.NewCurrencyAPIChain(
		currency_client.NewExchangeRateAPIClient(httpClient, env.ExchangeCurrencyAPIKey, logger),
	)
	currencyBeaconAPIClient := currency_client.NewCurrencyAPIChain(
		currency_client.NewCurrencyBeaconAPIClient(httpClient, env.CurrencyBeaconAPIKey, logger),
	)
	freeCurrencyAPIClient := currency_client.NewCurrencyAPIChain(
		currency_client.NewFreeCurrencyAPIClient(httpClient, logger),
	)
	exchangeRateAPIClient.SetNext(currencyBeaconAPIClient)
	currencyBeaconAPIClient.SetNext(freeCurrencyAPIClient)

	// initalization currency service server
	currencyService := currency_service.NewCurrencyService(exchangeRateAPIClient)
	currencyServiceServer := currency_service.NewCurrencyServiceServer(
		currencyService,
		logger,
	)

	// starting of server
	grpcServer := grpc.NewServer()
	currency_grpc.RegisterCurrencyServiceServer(grpcServer, currencyServiceServer)

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("currency service started at %s", url)
	utils.PanicOnError(grpcServer.Serve(lis), "failed to serve")
}
