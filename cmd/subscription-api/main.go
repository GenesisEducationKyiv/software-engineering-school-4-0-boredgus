package main

import (
	"flag"
	"strings"
	"subscription-api/cmd/subscription-api/internal"
	"subscription-api/config"
	pb_cs "subscription-api/pkg/grpc/currency_service"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	currencyServiceConn, err := grpc.NewClient(
		env.CurrencyServiceServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.FatalOnError(err, logger, "failed to connect to currency service grpc server")

	dispatchServiceConn, err := grpc.NewClient(
		env.DispatchServiceServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.FatalOnError(err, logger, "failed to connect to dispatch service grpc server")

	logger.Infof("started subscription API at %v port", env.Port)
	if err := internal.GetRouter(internal.APIParams{
		CurrencyService: pb_cs.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService: pb_ds.NewDispatchServiceClient(dispatchServiceConn),
		Logger:          logger,
	}).Run(":" + env.Port); err != nil {
		logger.Fatal("failed to start server: %v", err)
	}
}
