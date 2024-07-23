package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/currency"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/clients/mailman"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier"
	vm "github.com/VictoriaMetrics/metrics"
	"github.com/iancoleman/strcase"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	microserviceName    string        = "notification-service"
	metricsPushInterval time.Duration = 15 * time.Second
)

func main() {
	env, err := config.GetEnv()
	panicOnError(err, "failed to init environment variables")

	logger := config.InitLogger(env.Mode, config.WithProcess(microserviceName))
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
	emailNotifier := notifier.NewEmailNotifier(mailmanClient)
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

	// initialization of consumer
	consumer, err := broker.CreateNotificationConsumer(js)
	panicOnError(err, "failed to init consumer")

	// initialization of broker client
	natsBroker, err := broker.NewNatsBroker(js, consumer, logger)
	panicOnError(err, "failed to create broker")
	brokerWithMetrics := broker.NewNatsBrokerWithMetrics(natsBroker, broker.Metrics{
		TotalMessages:              vm.NewCounter(strcase.ToSnake("total consumed messages")),
		SuccessfulyHandledMessages: vm.NewCounter(strcase.ToSnake("successfuly handled messages")),
	})

	// connection to object store
	jetstreamStore, err := brokerWithMetrics.ObjectStore("dispatches")
	panicOnError(err, "failed to connect to object store")

	dispatchRepo := repo.NewDispatchRepo(broker.NewObjectStore(jetstreamStore))
	handler := app.NewEventHandler(
		brokerWithMetrics,
		notificationService,
		logger,
		dispatchRepo,
	)

	// initalization of dispatch scheduler
	scheduler := scheduler.NewDispatchScheduler(
		dispatchRepo,
		brokerWithMetrics,
		currencyService,
		logger,
	)
	defer scheduler.Stop()

	// initialization of metrics push
	commonLabels := map[string]string{"service": microserviceName}
	vm.ExposeMetadata(true)
	err = vm.InitPush(
		env.MetricsGatewayURL,
		metricsPushInterval,
		parseMetricLabels(commonLabels),
		true)
	panicOnError(err, "failed to init metrics push")

	// start the app
	app.NewApp(
		handler,
		scheduler,
		logger,
	).Run()
}

func panicOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err.Error()))
	}
}

func parseMetricLabels(labels map[string]string) string {
	labelsArr := make([]string, 0, len(labels))
	for key, value := range labels {
		labelsArr = append(labelsArr, fmt.Sprintf(`%s="%s"`, key, value))
	}

	return strings.Join(labelsArr, ",")
}
