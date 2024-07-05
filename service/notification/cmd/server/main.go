package main

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/currency"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/mailman"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier"
	"github.com/nats-io/nats.go"
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
	panicOnError(err, "failed to connect to dispatch service grpc server")
	defer currencyServiceConn.Close()

	// initialization of notification service
	mailmanClient := mailman.NewMailman(mailman.SMTPParams{
		Host:     env.SMTPHost,
		Port:     env.SMTPPort,
		Email:    env.SMTPEmail,
		Name:     env.SMTPUsername,
		Password: env.SMTPPassword,
	})
	baseNotifier := notifier.NewBaseNotifier()
	emailNotifier := notifier.NewEmailNotifier(baseNotifier, mailmanClient)

	notificationService := service.NewNotificationService(
		logger,
		emailNotifier,
		currency.NewCurrencyServiceClient(currencyServiceConn),
	)

	// connection to NATS broker
	natsConnection, err := nats.Connect(
		env.BrokerURL,
		nats.Name("notification-service"),
	)
	panicOnError(err, "failed to connect to NATS broker")

	natsBroker, err := broker.NewNatsBroker(natsConnection, logger)
	panicOnError(err, "failed to create broker")

	// initalzation of cron scheduler
	scheduler := scheduler.NewDispatchScheduler(func(d *entities.Dispatch) {
		// TODO: implement
	})

	jetstreamStore, err := natsBroker.ObjectStore("dispatches")
	panicOnError(err, "failed to connect to object store")

	dispatchStore := repo.NewDispatchRepo(broker.NewObjectStore(jetstreamStore))
	handler := app.NewEventHandler(
		natsBroker,
		dispatchStore,
		scheduler,
		notificationService,
		logger,
	)

	app.NewApp(
		handler,
		scheduler,
		logger,
		dispatchStore,
	).Run()
}
