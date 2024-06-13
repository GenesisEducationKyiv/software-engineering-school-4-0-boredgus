package main

import (
	"flag"
	"net"
	"strings"
	"subscription-api/cmd/dispatch-service/internal"
	"subscription-api/config"
	store "subscription-api/internal/db"
	ds "subscription-api/internal/services/dispatch"
	g "subscription-api/internal/services/dispatch/grpc"
	"subscription-api/internal/sql"
	"subscription-api/pkg/db"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
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

	lis, err := net.Listen("tcp", env.DispatchServiceServer)
	utils.FatalOnError(err, logger, "failed to listen: %v")

	server := grpc.NewServer()
	pb_ds.RegisterDispatchServiceServer(server,
		g.NewDispatchServiceServer(
			ds.NewDispatchService(
				store.NewStore(
					utils.Must(db.NewPostrgreSQL(
						env.PostgreSQLConnString,
						sql.PostgeSQLMigrationsUp("public"),
					)),
					db.IsPqError,
				))),
	)

	logger.Info("dispatch service started...")
	err = server.Serve(lis)
	utils.FatalOnError(err, logger, "failed to serve: %v")
}
