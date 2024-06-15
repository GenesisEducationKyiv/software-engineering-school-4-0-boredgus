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

	pb_cs "subscription-api/pkg/grpc/currency_service"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode)

	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to currency service grpc server")

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	server := grpc.NewServer()
	pb_ds.RegisterDispatchServiceServer(server,
		g.NewDispatchServiceServer(
			ds.NewDispatchService(
				ds.DispatchServiceParams{
					Store: store.NewStore(
						utils.Must(db.NewPostrgreSQL(
							env.PostgreSQLConnString,
							sql.PostgeSQLMigrationsUp("public", logger),
						)),
						db.IsPqError,
					),
					Logger: logger,
					Mailman: mailing.NewMailman(
						mailing.SMTPParams{
							Host:     env.SMTPHost,
							Port:     env.SMTPPort,
							Email:    env.SMTPEmail,
							Name:     env.SMTPUsername,
							Password: env.SMTPPassword,
						},
					),
					CurrencyService: pb_cs.NewCurrencyServiceClient(currencyServiceConn),
				},
			),
			logger,
		))

	logger.Infof("dispatch service started at %s", url)
	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
