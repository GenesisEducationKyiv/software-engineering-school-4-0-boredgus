package main

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/config"
	dispatch_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/pkg/grpc/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(config.Env())
	logger := config.InitLogger(env.Mode, config.WithProcess("dispatch-daemon"))
	defer logger.Flush()

	connURL := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	logger.Info(connURL)

	dispatchServiceConn, err := grpc.NewClient(
		connURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to dispatch service grpc server")
	defer dispatchServiceConn.Close()

	logger.Info("dispatch daemon has started")

	daemon := app.NewDispatchDaemon(
		dispatch_client.NewDispatchServiceClient(dispatchServiceConn),
		logger,
		app.NewScheduler(logger),
	)

	daemon.Run(context.Background())
}
