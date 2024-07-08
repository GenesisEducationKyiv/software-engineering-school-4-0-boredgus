package main

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/currency"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/mailman"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess("notification-service"))
	defer logger.Flush()

	// connection to gRPC server of currency service
	currencyServiceConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", env.CurrencyServiceAddress, env.CurrencyServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	panicOnError(err, "failed to connect to currency service grpc server")
	defer currencyServiceConn.Close()

	currencyService := currency.NewCurrencyServiceClient(currencyServiceConn)

	// initialization of notification service
	mailmanClient := mailman.NewSMTPMailman(mailman.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Name:     env.SMTPUsername,
		Password: env.SMTPPassword,
	})
	emailNotifier := notifier.NewEmailNotifier(
		notifier.NewBaseNotifier(),
		mailmanClient,
	)
	notificationService := service.NewNotificationService(logger, emailNotifier)

	// connection to NATS broker
	natsConnection, err := nats.Connect(
		env.BrokerURL,
		nats.Name("notification-service"),
	)
	panicOnError(err, "failed to connect to NATS broker")

	// initialization of jetstream
	js, err := jetstream.New(natsConnection)
	panicOnError(err, "failed to create NATS Jetstream instance")

	// initialization of broker client
	natsBroker, err := broker.NewNatsBroker(js, logger)
	panicOnError(err, "failed to create broker client")

	// connection to object store
	jetstreamStore, err := natsBroker.ObjectStore("dispatches")
	panicOnError(err, "failed to connect to object store")

	// initalization of dispatch scheduler
	scheduler := scheduler.NewDispatchScheduler(natsBroker, currencyService, logger)

	dispatchRepo := repo.NewDispatchRepo(broker.NewObjectStore(jetstreamStore))
	handler := app.NewEventHandler(
		natsBroker,
		scheduler,
		notificationService,
		logger,
		dispatchRepo,
	)

	app.NewApp(
		handler,
		scheduler,
		logger,
		dispatchRepo,
	).Run()
}
