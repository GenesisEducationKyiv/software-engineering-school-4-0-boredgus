package broker

import (
	"time"

	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
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
	CancelSubscriptionSubject string = "events.subscription.cancelled"

	CreateSubscriptionEvent string = "SubscriptionCreated"
	CancelSubscriptionEvent string = "SubscriptionCancelled"
)

func NewEventBroker(broker Broker, logger config.Logger) *eventBroker {
	return &eventBroker{
		broker: broker,
		logger: logger,
	}
}

func (b *eventBroker) CreateSubscription(sub service.Subscription) {
	event := CreateSubscriptionEvent

	data, err := proto.Marshal(&messages.SubscriptionMessage{
		EventType: messages.EventType_SUBSCRIPTION_CREATED,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   subscriptionToProto(sub, messages.SubscriptionStatus_CREATED),
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

func subscriptionToProto(
	sub service.Subscription,
	status messages.SubscriptionStatus,
) *messages.Subscription {
	return &messages.Subscription{
		DispatchID:  sub.DispatchID,
		BaseCcy:     sub.BaseCcy,
		TargetCcies: sub.TargetCcies,
		Email:       sub.Email,
		SendAt:      timestamppb.New(sub.SendAt),
		Status:      status,
	}
}
