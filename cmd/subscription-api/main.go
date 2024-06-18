package main

import (
	"fmt"
	"subscription-api/cmd/subscription-api/internal"
	"subscription-api/config"
	grpc_clients "subscription-api/pkg/grpc"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode, config.WithProcess("api"))
	defer logger.Flush()

	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	dispatchServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to dispatch service grpc server")
	defer dispatchServiceConn.Close()

	logger.Infof("started subscription API at %v port", env.Port)

	err = internal.GetRouter(&internal.APIParams{
		CurrencyService: grpc_clients.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService: grpc_clients.NewDispatchServiceClient(dispatchServiceConn),
		Logger:          logger,
	}).Run(":" + env.Port)
	utils.PanicOnError(err, "failed to start server")
}
