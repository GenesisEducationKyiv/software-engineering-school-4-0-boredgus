package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/dispatch-service/internal"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/mailing"
	dispatch_service "subscription-api/internal/services/dispatch"
	dispatch_grpc "subscription-api/internal/services/dispatch/grpc"
	"subscription-api/internal/sql"
	"subscription-api/pkg/db"

	grpc_client "subscription-api/pkg/grpc"
	"subscription-api/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(internal.Env())
	logger := config.InitLogger(env.Mode, config.WithProcess("dispatch-service"))
	defer logger.Flush()

	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	postgresqlDB, err := db.NewPostrgreSQL(
		env.PostgreSQLConnString,
		sql.PostgeSQLMigrationsUp(logger),
	)
	utils.PanicOnError(err, "failed toconnect to postgresql db")
	defer postgresqlDB.Close()

	store := store.NewStore(postgresqlDB, db.IsPqError)

	mailman := mailing.NewMailman(mailing.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Name:     env.SMTPUsername,
		Password: env.SMTPPassword,
	})

	serviceParams := &dispatch_service.DispatchServiceParams{
		Store:           store,
		Logger:          logger,
		Mailman:         mailman,
		CurrencyService: grpc_client.NewCurrencyServiceClient(currencyServiceConn),
	}

	dispatchServiceServer := dispatch_service.NewDispatchServiceServer(
		dispatch_service.NewDispatchService(serviceParams),
		logger,
	)

	server := grpc.NewServer()
	dispatch_grpc.RegisterDispatchServiceServer(server, dispatchServiceServer)

	logger.Infof("dispatch service started at %s", url)
	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
