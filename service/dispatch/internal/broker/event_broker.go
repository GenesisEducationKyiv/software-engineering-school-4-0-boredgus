package broker

import (
	"time"

	broker_msgs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	"github.com/google/uuid"
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

const (
	CreateSubscriptionSubject string = "events.subscription.created"

	CreateSubscriptionEvent string = "SubscriptionCreated"
)

func NewEventBroker(broker Broker, logger config.Logger) *eventBroker {
	return &eventBroker{
		broker: broker,
		logger: logger,
	}
}

func (b *eventBroker) CreateSubscription(sub deps.Subscription) {
	event := CreateSubscriptionEvent

	data, err := proto.Marshal(&broker_msgs.SubscriptionCreatedMessage{
		EventID:   uuid.NewString(),
		EventType: event,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   SubscriptionToProto(sub, broker_msgs.SubscriptionStatus_CREATED),
	})
	if err != nil {
		b.logger.Errorf("failed to marshal %s message: %v", event, err)

		return
	}

	if err = b.broker.PublishAsync(CreateSubscriptionSubject, data); err != nil {
		b.logger.Errorf(
			"failed to publish CreateSubscription message to '%s' subject: %v",
			CreateSubscriptionSubject, err,
		)
	}
}

func SubscriptionToProto(
	sub deps.Subscription,
	status broker_msgs.SubscriptionStatus,
) *broker_msgs.SubscriptionPayload {
	return &broker_msgs.SubscriptionPayload{
		DispatchID:  sub.DispatchID,
		BaseCcy:     sub.BaseCcy,
		TargetCcies: sub.TargetCcies,
		Email:       sub.Email,
		SendAt:      timestamppb.New(sub.SendAt),
		Status:      status,
	}
}
