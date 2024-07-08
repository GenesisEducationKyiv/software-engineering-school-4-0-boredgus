package integration

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/tests"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/tests/stubs"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/suite"
)

const (
	ConsumerName string = "test-consumer"
)

type (
	MessageHandler interface {
		HandleMessages() error
	}

	NATSConsumerSuite struct {
		*suite.Suite

		natsContainer *tests.NATSContainer
		ctx           context.Context
		logger        config.Logger

		// handler MessageHandler
		broker app.Consumer
	}
)

func (t *NATSConsumerSuite) initConsumerClient() {
	t.logger = config.InitLogger(config.TestMode)
	natsConnection, err := nats.Connect(
		t.natsContainer.ConnectionString,
		nats.Name(ConsumerName),
	)
	t.NoErrorf(err, "failed to connect to NATS broker")

	js, err := jetstream.New(natsConnection)
	t.NoErrorf(err, "failed to create NATS Jetstream instance")

	natsBroker, err := broker.NewNatsBroker(js, t.logger)
	t.NoErrorf(err, "failed to create broker client")
	t.broker = natsBroker
}

func (t *NATSConsumerSuite) SetupSuite() {
	t.ctx = context.Background()
	natsContainer, err := tests.CreateNatsContainer(t.ctx)
	t.NoErrorf(err, "failed to create NATS container")

	t.natsContainer = natsContainer
	t.initConsumerClient()

	handler := app.NewEventHandler(
		t.broker,
		stubs.NewSchedulerMock(),
		service.NewNotificationService(t.logger, stubs.NewNotifierMock()),
		t.logger,
		stubs.NewDispatchStoreMock(),
	)
	t.NoErrorf(handler.HandleMessages(), "failed to handle messages")

}
func (t *NATSConsumerSuite) AfterTest() {

}

func (t *NATSConsumerSuite) TearDownSuite() {
	t.NoError(t.natsContainer.Container.Terminate(t.ctx))
}

func (t *NATSConsumerSuite) TestConsumingOf_SubscriptionCreatedEvent() {

}

func TestNATSConsumerSuite(t *testing.T) {
	suite.Run(t, new(NATSConsumerSuite))
}
