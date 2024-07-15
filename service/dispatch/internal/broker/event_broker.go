package broker

import (
	"reflect"
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

var statusToSubject = map[service.SubscriptionStatus]string{
	service.SubscriptionStatusActive:    "events.subscription.created",
	service.SubscriptionStatusCancelled: "events.subscription.cancelled",
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

	status := service.SubscriptionStatus(statusValue.Int())
	subject, ok := statusToSubject[status]
	if !ok {
		b.logger.Error("provided status is unsupported")

		return
	}

	data, err := proto.Marshal(&messages.SubscriptionMessage{
		EventType: statusToEventType[status],
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload:   subscriptionToProto(msg.(service.Subscription), statusToProtoStatus[status]),
	})
	if err != nil {
		b.logger.Errorf("failed to marshal subscription message: %v", err)

		return
	}

	if err = b.broker.PublishAsync(subject, data); err != nil {
		b.logger.Errorf("failed to publish message to '%s' subject: %v", subject, err)
	}
}

var statusToEventType = map[service.SubscriptionStatus]messages.EventType{
	service.SubscriptionStatusActive:    messages.EventType_SUBSCRIPTION_CREATED,
	service.SubscriptionStatusCancelled: messages.EventType_SUBSCRIPTION_CANCELLED,
}

var statusToProtoStatus = map[service.SubscriptionStatus]messages.SubscriptionStatus{
	service.SubscriptionStatusActive:    messages.SubscriptionStatus_CREATED,
	service.SubscriptionStatusCancelled: messages.SubscriptionStatus_CANCELLED,
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
