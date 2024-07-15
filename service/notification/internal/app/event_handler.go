package app

import (
	"context"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	broker_msgs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"google.golang.org/protobuf/proto"
)

type (
	Consumer interface {
		ConsumeMessage(handler func(msg broker.ConsumedMessage) error) error
	}

	NotificationService interface {
		SendSubscriptionDetails(ctx context.Context, notification service.Notification) error
		SendExchangeRates(ctx context.Context, notification service.Notification) error
	}

	DispatchStore interface {
		AddSubscription(ctx context.Context, sub *entities.Subscription) error
		CancelSubscription(ctx context.Context, sub *entities.Subscription) error
	}

	Scheduler interface {
		AddSubscriberToDispatch(*entities.Subscription)
		RemoveSubscriberFromDispatch(email, dispatchID string)
	}

	eventHandler struct {
		broker        Consumer
		scheduler     Scheduler
		logger        config.Logger
		service       NotificationService
		dispatchStore DispatchStore
	}
)

const (
	TimeoutOfProcessing time.Duration = 2 * time.Second
)

func NewEventHandler(
	broker Consumer,
	dispatchScheduler Scheduler,
	service NotificationService,
	logger config.Logger,
	dispatchStore DispatchStore,
) *eventHandler {

	return &eventHandler{
		logger:        logger,
		scheduler:     dispatchScheduler,
		broker:        broker,
		service:       service,
		dispatchStore: dispatchStore,
	}
}

func (h *eventHandler) HandleMessages() error {
	return h.broker.ConsumeMessage(func(msg broker.ConsumedMessage) error {
		return h.handlerFactory(msg.Subject())(msg)
	})
}

func (h *eventHandler) handlerFactory(subject string) func(broker.ConsumedMessage) error {
	switch subject {
	case SubscriptionCreatedEvent, SubscriptionCancelledEvent:
		return h.handleSubscriptionEvent
	case SendDispatchCommand:
		return h.handleSendDispatchCommand
	}

	return func(cm broker.ConsumedMessage) error {
		return broker.ErrSkippedMessage
	}
}

func (h *eventHandler) handleSubscriptionEvent(msg broker.ConsumedMessage) error {
	var parsedMsg broker_msgs.SubscriptionMessage
	if err := proto.Unmarshal(msg.Data(), &parsedMsg); err != nil {
		return fmt.Errorf("failed to unmarshal subscription message: %w", err)
	}

	sub := ProtoToSubscription(&parsedMsg)

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	var err error
	switch parsedMsg.EventType {
	case broker_msgs.EventType_SUBSCRIPTION_CREATED:
		h.scheduler.AddSubscriberToDispatch(sub)
		err = h.dispatchStore.AddSubscription(ctx, sub)
	case broker_msgs.EventType_SUBSCRIPTION_CANCELLED:
		h.scheduler.RemoveSubscriberFromDispatch(sub.Email, sub.DispatchID)
		err = h.dispatchStore.CancelSubscription(ctx, sub)
	}

	if err != nil {
		return fmt.Errorf("failed to save dispatches: %w", err)
	}

	if err := h.service.SendSubscriptionDetails(
		ctx,
		*ProtoToDispatchDetailsNotification(&parsedMsg, MessageTypeToNotificationType(parsedMsg.EventType)),
	); err != nil {
		return fmt.Errorf("failed to send subscription details: %w", err)
	}

	return nil
}

func (h *eventHandler) handleSendDispatchCommand(msg broker.ConsumedMessage) error {
	var parsedMsg broker_msgs.SendDispatchCommand
	if err := proto.Unmarshal(msg.Data(), &parsedMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message from %s: %w", SubscriptionCreatedEvent, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	if err := h.service.SendExchangeRates(
		ctx,
		*ProtoToCurrencyDispatchNotification(&parsedMsg),
	); err != nil {
		return fmt.Errorf("failed to send subscription details: %w", err)
	}

	return nil
}
