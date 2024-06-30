package main

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
	currency_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/pkg/grpc/client"
	dispatch_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/pkg/grpc/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(config.Env())
	logger := logger.InitLogger(env.Mode, logger.WithProcess("api"))
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

	router := config.GetRouter(&config.APIParams{
		CurrencyService: currency_client.NewCurrencyServiceClient(currencyServiceConn),
		DispatchService: dispatch_client.NewDispatchServiceClient(dispatchServiceConn),
		Logger:          logger,
	})

	utils.PanicOnError(router.Run(":"+env.Port), "failed to start server")
}
