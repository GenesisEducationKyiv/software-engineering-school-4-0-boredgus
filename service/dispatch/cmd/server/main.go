package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/nats-io/nats.go"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/server"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"google.golang.org/grpc"
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

	// connection to db
	postgresqlDB, err := db.NewPostrgreSQL(
		env.PostgreSQLConnString,
		db.PostgeSQLMigrationsUp,
	)
	panicOnError(err, "failed to connect to postgresql db")
	defer postgresqlDB.Close()

	storage := repo.NewStore(postgresqlDB, db.IsPqError)

	// connection to NATS broker
	natsConnection, err := nats.Connect(
		env.BrokerURL,
		nats.Name("subscription-service"),
	)
	panicOnError(err, "failed to connect to broker")

	natsBroker := broker.NewNatsBroker(natsConnection, logger, panicOnError)

	// initialization of dispatch service server
	dispatchService := service.NewDispatchService(
		repo.NewUserRepo(storage),
		repo.NewSubRepo(storage),
		repo.NewDispatchRepo(storage),
		broker.NewEventBroker(natsBroker, logger),
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
