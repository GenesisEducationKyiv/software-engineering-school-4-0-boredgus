package main

import (
	"fmt"
	"net"
	"subscription-api/cmd/dispatch-service/internal"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/mailing"
	dispatch_service "subscription-api/internal/services/dispatch"
	dispatch_server "subscription-api/internal/services/dispatch/server"
	dispatch_grpc "subscription-api/internal/services/dispatch/server/grpc"
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

	// connection to currency service server
	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.PanicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	// connection to db
	postgresqlDB := utils.Must(db.NewPostrgreSQL(
		env.PostgreSQLConnString,
		sql.PostgeSQLMigrationsUp(logger),
	))

	utils.PanicOnError(err, "failed toconnect to postgresql db")
	defer postgresqlDB.Close()

	storage := store.NewStore(postgresqlDB, db.IsPqError)

	smtpParams := mailing.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Name:     env.SMTPUsername,
		Password: env.SMTPPassword,
	}

	// initialization of dispatch service server
	dispatchService := dispatch_service.NewDispatchService(
		logger,
		mailing.NewMailman(smtpParams),
		grpc_client.NewCurrencyServiceClient(currencyServiceConn),
		store.NewUserRepo(storage),
		store.NewSubRepo(storage),
		store.NewDispatchRepo(storage),
	)
	dispatchServiceServer := dispatch_server.NewDispatchServiceServer(dispatchService, logger)

	// starting of grpc server
	server := grpc.NewServer()
	dispatch_grpc.RegisterDispatchServiceServer(server, dispatchServiceServer)

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("dispatch service started at %s", url)
	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
