package integration

import (
	"context"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/tests"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/tests/stubs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/tests/testdata"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

const (
	ConsumerName string = "test-consumer"
)

type (
	MessageHandler interface {
		HandleMessages() error
	}
	Broker interface {
		ConsumeMessage(handler func(msg broker.ConsumedMessage) error) error
		PublishAsync(subject string, payload []byte) error
	}

	NATSConsumerSuite struct {
		suite.Suite

		natsContainer *tests.NATSContainer
		ctx           context.Context
		logger        config.Logger

		broker  Broker
		service *stubs.NotificationServiceMock
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

	consumer, err := broker.CreateNotificationConsumer(js)
	t.NoErrorf(err, "failed to create consumer")

	natsBroker, err := broker.NewNatsBroker(js, consumer, t.logger)
	t.NoErrorf(err, "failed to create broker client")
	t.broker = natsBroker
}

func (t *NATSConsumerSuite) marshalMessage(m proto.Message) []byte {
	data, err := proto.Marshal(m)
	t.NoErrorf(err, "failed to marshal message")

	return data
}

func (t *NATSConsumerSuite) SetupSuite() {
	t.ctx = context.Background()
	natsContainer, err := tests.CreateNatsContainer(t.ctx)
	t.NoErrorf(err, "failed to create NATS container")

	t.natsContainer = natsContainer
	t.initConsumerClient()

	t.service = stubs.NewNotificationServiceMock()

	handler := app.NewEventHandler(
		t.broker,
		stubs.NewSchedulerMock(),
		t.service,
		t.logger,
		stubs.NewDispatchStoreMock(),
	)
	t.NoErrorf(handler.HandleMessages(), "failed to handle messages")
}

func (t *NATSConsumerSuite) TearDownSuite() {
	t.NoError(t.natsContainer.Container.Terminate(t.ctx))
}

func (t *NATSConsumerSuite) TestConsumingOf_SubscriptionCreatedEvent() {
	msg := testdata.SubscriptionCreatedMessage
	err := t.broker.PublishAsync(app.SubscriptionCreatedEvent, t.marshalMessage(msg))

	t.NoErrorf(err, "failed to publish message to broker")
	time.Sleep(300 * time.Millisecond)
	t.service.AssertNotificationTypeOfLastCall(t.T(), service.SubscriptionCreated)
}

func (t *NATSConsumerSuite) TestConsumingOf_SubscriptionCancelledEvent() {
	msg := testdata.SubscriptionCancelledMessage
	err := t.broker.PublishAsync(app.SubscriptionCancelledEvent, t.marshalMessage(msg))

	t.NoErrorf(err, "failed to publish message to broker")
	time.Sleep(300 * time.Millisecond)
	t.service.AssertNotificationTypeOfLastCall(t.T(), service.SubscriptionCancelled)
}

func TestNATSConsumerSuite(t *testing.T) {
	suite.Run(t, new(NATSConsumerSuite))
}
