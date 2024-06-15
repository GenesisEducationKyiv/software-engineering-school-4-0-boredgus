package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/dispatch-service/internal"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/mailing"
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
	smtpParams := mailing.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Password: env.SMTPPassword,
	}
	pb_ds.RegisterDispatchServiceServer(server,
		g.NewDispatchServiceServer(
			ds.NewDispatchService(
				store.NewStore(
					utils.Must(db.NewPostrgreSQL(
						env.PostgreSQLConnString,
						sql.PostgeSQLMigrationsUp("public", logger),
					)),
					db.IsPqError,
				),
				logger,
				smtpParams),
			logger,
		))

	logger.Info("dispatch service started...")
	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
