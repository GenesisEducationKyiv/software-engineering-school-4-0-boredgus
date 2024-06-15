package main

import (
	"context"
	"fmt"
	"subscription-api/cmd/dispatch-daemon/internal"
	"subscription-api/config"
	"subscription-api/pkg/utils"

	pb_ds "subscription-api/pkg/grpc/dispatch_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode).With("service", "dispatch-daemon")

	connURL := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	logger.Info(connURL)

	dispatchServiceConn, err := grpc.NewClient(
		connURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to dispatch service grpc server")
	logger.Info(dispatchServiceConn)

	logger.Info("dispatch daemon has started")
	internal.NewDispatchDaemon(
		pb_ds.NewDispatchServiceClient(dispatchServiceConn),
		logger,
		internal.NewScheduler(logger),
	).Run(context.Background())
}
