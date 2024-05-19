package main

// d := mail.NewDialer("smtp.gmail.com", 587, "daha.kyiv@gmail.com", "guze dokh umzh ulvs")

import (
	"flag"
	"fmt"
	"os"
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
	config.InitEnvVariables(strings.Split(envFile, ",")...)
}

func main() {
	logger := config.InitLogger(config.Mode(os.Getenv("MODE")))
	logger.Info("START API at ", os.Getenv("API_PORT"))

	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%v:%v", os.Getenv("CURRENCY_SERVICE_ADDRESS"), os.Getenv("CURRENCY_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.FatalOnError(err, logger, "failed to connect to currency service grpc server")

	dispatchServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%v:%v", os.Getenv("DISPATCH_SERVICE_ADDRESS"), os.Getenv("DISPATCH_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.FatalOnError(err, logger, "failed to connect to dispatch service grpc server")

	if err := internal.GetRouter(internal.APIParams{
		CurrencyService: pb_cs.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService: pb_ds.NewDispatchServiceClient(dispatchServiceConn),
		Logger:          logger,
	}).Run(":" + os.Getenv("API_PORT")); err != nil {
		logger.Fatal("failed to start server: %v", err)
	}
}
