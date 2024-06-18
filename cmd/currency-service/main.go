package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/currency-service/internal"
	"subscription-api/config"
	"subscription-api/internal/clients"
	currency_service "subscription-api/internal/services/currency"
	currency_grpc "subscription-api/internal/services/currency/grpc"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode, config.WithProcess("currency-service"))
	defer logger.Flush()

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	currencyService := currency_service.NewCurrencyService(
		clients.NewExchangeRateAPIClient(env.ExchangeCurrencyAPIKey),
	)

	currencyServiceServer := currency_service.NewCurrencyServiceServer(
		currencyService,
		logger,
	)

	grpcServer := grpc.NewServer()
	currency_grpc.RegisterCurrencyServiceServer(grpcServer, currencyServiceServer)
	logger.Info("currency service started...")

	utils.PanicOnError(grpcServer.Serve(lis), "failed to serve")
}
