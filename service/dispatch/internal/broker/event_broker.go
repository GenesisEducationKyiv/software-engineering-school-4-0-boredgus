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
	SubscriptionCreatedSubject   string = "events.subscription.created"
	SubscriptionCancelledSubject string = "events.subscription.cancelled"
)

func NewEventBroker(broker Broker, logger config.Logger) *eventBroker {
	return &eventBroker{
		broker: broker,
		logger: logger,
	}
}

func (b *eventBroker) CreateSubscription(sub service.Subscription) {
	data, err := proto.Marshal(&messages.SubscriptionMessage{
		EventType: messages.EventType_SUBSCRIPTION_CREATED,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   subscriptionToProto(sub, messages.SubscriptionStatus_CREATED),
	})
	if err != nil {
		b.logger.Errorf("failed to marshal subscription message: %v", err)

		return
	}

	subject := SubscriptionCreatedSubject
	if err = b.broker.PublishAsync(subject, data); err != nil {
		b.logger.Errorf("failed to publish message to '%s' subject: %v", subject, err)
	}
}

func (b *eventBroker) CancelSubscription(sub service.Subscription) {
	data, err := proto.Marshal(&messages.SubscriptionMessage{
		EventType: messages.EventType_SUBSCRIPTION_CANCELLED,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   subscriptionToProto(sub, messages.SubscriptionStatus_CANCELLED),
	})
	if err != nil {
		b.logger.Errorf("failed to marshal subscription message: %v", err)

		return
	}

	subject := SubscriptionCancelledSubject
	if err = b.broker.PublishAsync(subject, data); err != nil {
		b.logger.Errorf("failed to publish message to '%s' subject: %v", subject, err)
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
