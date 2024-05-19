package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
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

var envFile string

func init() {
	flag.StringVar(&envFile, "env", "dev.env", "list of filenames splitted with coma (e.g. '.env,dev.env')")
	flag.Parse()
	config.InitEnvVariables(strings.Split(envFile, ",")...)
}

func main() {
	logger := config.InitLogger(config.DevMode)

	address := fmt.Sprintf("%v:%v", os.Getenv("DISPATCH_SERVICE_ADDRESS"), os.Getenv("DISPATCH_SERVICE_PORT"))
	lis, err := net.Listen("tcp", address)
	utils.FatalOnError(err, logger, "failed to listen: %v")

	dbName := os.Getenv("POSTGRESQL_DB")
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s search_path=%s sslmode=disable",
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_ADDRESS"),
		dbName, dbName)

	server := grpc.NewServer()
	pb_ds.RegisterDispatchServiceServer(server,
		g.NewDispatchServiceServer(
			ds.NewDispatchService(d_store.NewCurrencyDispatchStore(
				store.NewStore(utils.Must(db.NewPostrgreSQL(dsn, sql.PostgeSQLMigrationsUp("public"))),
					db.IsPqError,
				))),
		))

	err = server.Serve(lis)
	utils.FatalOnError(err, logger, "failed to serve: %v")
}
