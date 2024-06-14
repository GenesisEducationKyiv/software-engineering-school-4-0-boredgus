package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/dispatch-service/internal"
	"subscription-api/config"
	store "subscription-api/internal/db"
	d_store "subscription-api/internal/db/dispatch"
	ds "subscription-api/internal/services/dispatch"
	g "subscription-api/internal/services/dispatch/grpc"
	"subscription-api/internal/sql"
	"subscription-api/pkg/db"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode)

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	server := grpc.NewServer()
	pb_ds.RegisterDispatchServiceServer(server,
		g.NewDispatchServiceServer(
			ds.NewDispatchService(d_store.NewCurrencyDispatchStore(
				store.NewStore(
					utils.Must(db.NewPostrgreSQL(
						env.PostgreSQLConnString,
						sql.PostgeSQLMigrationsUp("public", logger),
					)),
					db.IsPqError,
				))),
			logger,
		))

	logger.Info("dispatch service started...")

	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
