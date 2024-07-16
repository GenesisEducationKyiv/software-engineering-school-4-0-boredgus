package broker

import (
	"reflect"
	"time"

	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Broker interface {
		PublishAsync(subject string, data []byte) error
	}

	eventBroker struct {
		broker Broker
		logger config.Logger
	}
)

var statusToSubject = map[messages.SubscriptionStatus]string{
	messages.SubscriptionStatus_CREATED:   "events.subscription.created",
	messages.SubscriptionStatus_CANCELLED: "events.subscription.cancelled",
}

func NewEventBroker(broker Broker, logger config.Logger) *eventBroker {
	return &eventBroker{
		broker: broker,
		logger: logger,
	}
}

func (b *eventBroker) Publish(msg interface{}) {
	statusValue := reflect.ValueOf(msg).FieldByName("Status")
	if !statusValue.CanInt() {
		b.logger.Error("unsupported message provided")

		return
	}

	status := messages.SubscriptionStatus(statusValue.Int())
	subject, ok := statusToSubject[status]
	if !ok {
		b.logger.Error("provided status is unsupported")

		return
	}

	subscription := msg.(messages.Subscription)

	data, err := proto.Marshal(&messages.SubscriptionMessage{
		EventType: statusToEventType[status],
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   &subscription,
	})
	if err != nil {
		b.logger.Errorf("failed to marshal subscription message: %v", err)

		return
	}

	if err = b.broker.PublishAsync(subject, data); err != nil {
		b.logger.Errorf("failed to publish message to '%s' subject: %v", subject, err)
	}
}

var statusToEventType = map[messages.SubscriptionStatus]messages.EventType{
	messages.SubscriptionStatus_CREATED:   messages.EventType_SUBSCRIPTION_CREATED,
	messages.SubscriptionStatus_CANCELLED: messages.EventType_SUBSCRIPTION_CANCELLED,
}
