package main

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/db"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	"google.golang.org/grpc"
)

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess("customer-service"))

	database, err := db.NewDatabase(env.DatabaseURL)
	panicOnError(err, "failed to setup database")

	customerService := service.NewCustomerService(
		repo.NewCustomerRepo(database, db.IsError),
	)
	customerServer := server.NewCustomerServiceServer(customerService, logger)

	grpcServer := grpc.NewServer()
	grpc_gen.RegisterCustomerServiceServer(grpcServer, customerServer)

	url := fmt.Sprintf("%s:%s", env.CustomerServiceAddress, env.CustomerServicePort)
	lis, err := net.Listen("tcp", url)
	panicOnError(err, fmt.Sprintf("failed to listen %s", url))

	logger.Infof("customer service started at %s", url)
	panicOnError(grpcServer.Serve(lis), "failed to serve")
}
