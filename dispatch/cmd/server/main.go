package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/server"

	currency_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/pkg/grpc/client"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mailing"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/sql"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/db/postgres"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := utils.Must(config.Env())
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
	postgresqlDB := utils.Must(postgres.NewPostrgreSQL(
		env.PostgreSQLConnString,
		sql.PostgeSQLMigrationsUp(logger),
	))

	utils.PanicOnError(err, "failed toconnect to postgresql db")
	defer postgresqlDB.Close()

	storage := repo.NewStore(postgresqlDB, postgres.IsPqError)

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
	utils.PanicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("dispatch service started at %s", url)
	utils.PanicOnError(server.Serve(lis), "failed to serve")
}
