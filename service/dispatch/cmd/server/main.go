package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/server"

	currency_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/clients/currency"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mailing"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
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
	logger := config.InitLogger(env.Mode, config.WithProcess("dispatch-service"))
	defer logger.Flush()

	// connection to currency service server
	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	// connection to db
	postgresqlDB, err := db.NewPostrgreSQL(
		env.PostgreSQLConnString,
		db.PostgeSQLMigrationsUp(logger),
	)
	panicOnError(err, "failed toconnect to postgresql db")
	defer postgresqlDB.Close()

	storage := repo.NewStore(postgresqlDB, db.IsPqError)

	smtpParams := mailing.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Name:     env.SMTPUsername,
		Password: env.SMTPPassword,
	}

	// initialization of dispatch service server
	dispatchService := service.NewDispatchService(
		logger,
		mailing.NewMailman(smtpParams),
		currency_client.NewCurrencyServiceClient(currencyServiceConn),
		repo.NewUserRepo(storage),
		repo.NewSubRepo(storage),
		repo.NewDispatchRepo(storage),
	)
	dispatchServiceServer := server.NewDispatchServiceServer(dispatchService, logger)

	// starting of grpc server
	server := grpc.NewServer()
	grpc_gen.RegisterDispatchServiceServer(server, dispatchServiceServer)

	url := fmt.Sprintf("%s:%s", env.DispatchServiceAddress, env.DispatchServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("dispatch service started at %s", url)
	panicOnError(server.Serve(lis), "failed to serve")
}
