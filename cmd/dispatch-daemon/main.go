package main

import (
	"context"
	"fmt"
	"subscription-api/cmd/dispatch-daemon/internal"
	"subscription-api/config"
	"subscription-api/pkg/utils"

	grpc_client "subscription-api/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(internal.Env())
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
	internal.NewDispatchDaemon(
		grpc_client.NewDispatchServiceClient(dispatchServiceConn),
		logger,
		internal.NewScheduler(logger),
	).Run(context.Background())
}
