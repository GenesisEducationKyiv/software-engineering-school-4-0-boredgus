package main

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/clients/dispatch"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/config"
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
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess("dispatch-daemon"))
	defer logger.Flush()

	connURL := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	logger.Info(connURL)

	dispatchServiceConn, err := grpc.NewClient(
		connURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to dispatch service grpc server")
	defer dispatchServiceConn.Close()

	logger.Info("dispatch daemon has started")

	daemon := app.NewDispatchDaemon(
		dispatch.NewDispatchServiceClient(dispatchServiceConn),
		logger,
		app.NewScheduler(logger),
	)

	daemon.Run(context.Background())
}
