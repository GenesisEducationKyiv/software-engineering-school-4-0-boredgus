package main

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/currency"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/clients/dispatch"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
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

	logger := logger.InitLogger(env.Mode, logger.WithProcess("api"))
	defer logger.Flush()

	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	dispatchServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to dispatch service grpc server")
	defer dispatchServiceConn.Close()

	logger.Infof("started subscription API at %v port", env.Port)

	router := config.GetRouter(&config.APIParams{
		CurrencyService: currency.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService: dispatch.NewTransactionManagerClient(dispatchServiceConn),
		Logger:          logger,
	})

	panicOnError(router.Run(":"+env.Port), "failed to start server")
}
