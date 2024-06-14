package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/currency-service/internal"
	"subscription-api/config"
	cs "subscription-api/internal/services/currency"
	g "subscription-api/internal/services/currency/grpc"
	pb_cs "subscription-api/pkg/grpc/currency_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode)

	url := fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	server := grpc.NewServer()
	pb_cs.RegisterCurrencyServiceServer(server,
		g.NewCurrencyServiceServer(
			cs.NewCurrencyService(env.ExchangeCurrencyAPIKey),
			logger,
		))
	logger.Info("currency service started...")

	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
