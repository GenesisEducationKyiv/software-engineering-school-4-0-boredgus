package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"subscription-api/config"
	cs "subscription-api/internal/services/currency"
	g "subscription-api/internal/services/currency/grpc"
	pb_cs "subscription-api/pkg/grpc/currency_service"
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
	logger := config.InitLogger(config.Mode(os.Getenv("MODE")))

	address := fmt.Sprintf("%v:%v", os.Getenv("CURRENCY_SERVICE_ADDRESS"), os.Getenv("CURRENCY_SERVICE_PORT"))
	lis, err := net.Listen("tcp", address)
	utils.FatalOnError(err, logger, "failed to listen: %v")

	server := grpc.NewServer()
	pb_cs.RegisterCurrencyServiceServer(server,
		g.NewCurrencyServiceServer(
			cs.NewCurrencyService(os.Getenv("EXCHANGE_CURRENCY_API_KEY")),
		))
	err = server.Serve(lis)
	utils.FatalOnError(err, logger, "failed to serve: %v")
}
