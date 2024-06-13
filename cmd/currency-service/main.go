package main

import (
	"flag"
	"net"
	"strings"
	"subscription-api/cmd/currency-service/internal"
	"subscription-api/config"
	cs "subscription-api/internal/services/currency"
	g "subscription-api/internal/services/currency/grpc"
	pb_cs "subscription-api/pkg/grpc/currency_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
)

var envFile string

func init() {
	flag.StringVar(&envFile, "env", "dev.env", "list of filenames splitted with coma (e.g. '.env,dev.env')")
	flag.Parse()
	config.LoadEnvFile(strings.Split(envFile, ",")...)
}

func main() {
	env := internal.Env()
	logger := config.InitLogger(env.Mode)

	lis, err := net.Listen("tcp", env.CurrencyServiceServer)
	utils.FatalOnError(err, logger, "failed to listen: %v")

	server := grpc.NewServer()
	pb_cs.RegisterCurrencyServiceServer(server,
		g.NewCurrencyServiceServer(
			cs.NewCurrencyService(env.ExchangeCurrencyAPIKey),
		))
	logger.Info("currency service started...")
	err = server.Serve(lis)
	utils.FatalOnError(err, logger, "failed to serve: %v")
}
