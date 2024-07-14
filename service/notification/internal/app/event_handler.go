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
	}

	Scheduler interface {
		AddSubscription(*entities.Subscription)
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
	SubscriptionCreatedEvent string = "events.subscription.created"
	SendDispatchCommand      string = "commands.send.dispatch"

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
		logger:    logger,
		scheduler: dispatchScheduler,
		broker:    broker,
		service:   service,
	}
}

func (h *eventHandler) HandleMessages() error {
	return h.broker.ConsumeMessage(func(msg broker.ConsumedMessage) error {
		return h.handlerFactory(msg.Subject())(msg)
	})
}

func (h *eventHandler) handlerFactory(subject string) func(broker.ConsumedMessage) error {
	switch subject {
	case SubscriptionCreatedEvent:
		return h.handleSubscriptionCreatedEvent
	case SendDispatchCommand:
		return h.handleSendDispatchCommand
	}

	return func(cm broker.ConsumedMessage) error {
		return broker.SkippedMessageErr
	}
}

func (h *eventHandler) handleSubscriptionCreatedEvent(msg broker.ConsumedMessage) error {
	var parsedMsg broker_msgs.SubscriptionMessage
	if err := proto.Unmarshal(msg.Data(), &parsedMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message from %s: %w", SubscriptionCreatedEvent, err)
	}

	sub := ProtoToSubscription(&parsedMsg)

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	h.scheduler.AddSubscription(sub)
	if err := h.dispatchStore.AddSubscription(ctx, sub); err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}

	if err := h.service.SendSubscriptionDetails(
		ctx,
		*ProtoToDispatchDetailsNotification(&parsedMsg, service.SubscriptionCreated),
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
