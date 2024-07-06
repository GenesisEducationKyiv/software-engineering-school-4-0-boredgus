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
	Broker interface {
		PublishAsync(subject string, payload []byte) error
		ConsumeEvent(handler func(msg broker.ConsumedMessage)) error
		ConsumeCommand(handler func(msg broker.ConsumedMessage)) error
	}

	NotificationService interface {
		SendSubscriptionDetails(ctx context.Context, notification service.Notification) error
		SendExchangeRates(ctx context.Context, notification service.Notification) error
	}

	DispatchStore interface {
		AddSubscription(ctx context.Context, sub entities.Subscription) error
	}

	Scheduler interface {
		AddSubscription(entities.Subscription)
	}

	eventHandler struct {
		broker        Broker
		dispatchStore DispatchStore
		scheduler     Scheduler
		logger        config.Logger
		service       NotificationService
	}
)

const (
	SubscriptionCreatedEvent string = "events.subscription.created"
	SendDispatchCommand      string = "commands.send.dispatch"

	TimeoutOfProcessing time.Duration = 2 * time.Second
)

func NewEventHandler(
	broker Broker,
	dispatchStore DispatchStore,
	dispatchScheduler Scheduler,
	service NotificationService,
	logger config.Logger,
) *eventHandler {

	return &eventHandler{
		logger:        logger,
		dispatchStore: dispatchStore,
		scheduler:     dispatchScheduler,
		broker:        broker,
		service:       service,
	}
}

func (h *eventHandler) handleSubscriptionCreatedEvent(msg broker.ConsumedMessage) error {
	var parsedMsg broker_msgs.SubscriptionCreatedMessage
	if err := proto.Unmarshal(msg.Data(), &parsedMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message from %s: %w", SubscriptionCreatedEvent, err)
	}

	h.scheduler.AddSubscription(entities.Subscription{
		DispatchID:  parsedMsg.Payload.DispatchID,
		BaseCcy:     parsedMsg.Payload.BaseCcy,
		TargetCcies: parsedMsg.Payload.TargetCcies,
		Email:       parsedMsg.Payload.Email,
		SendAt:      parsedMsg.Payload.SendAt.AsTime(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	if err := h.service.SendSubscriptionDetails(ctx, service.Notification{
		Type: service.SubscriptionCreated,
		Data: service.NotificationData{
			Emails: []string{parsedMsg.Payload.Email},
			Payload: service.SubscriptionDetails{
				BaseCcy:     parsedMsg.Payload.BaseCcy,
				TargetCcies: parsedMsg.Payload.TargetCcies,
				SendAt:      parsedMsg.Payload.SendAt.AsTime().UTC().Format(time.TimeOnly),
			},
		},
	}); err != nil {
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

	if err := h.service.SendSubscriptionDetails(ctx, service.Notification{
		Type: service.SubscriptionCreated,
		Data: service.NotificationData{
			Emails: parsedMsg.Data.Emails,
			Payload: service.CurrencyDispatch{
				BaseCcy: parsedMsg.Data.BaseCcy,
				Rates:   parsedMsg.Data.Rates,
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send subscription details: %w", err)
	}

	return nil
}

func (h *eventHandler) HandleEvents() error {
	return h.broker.ConsumeEvent(func(msg broker.ConsumedMessage) {
		var err error

		switch msg.Subject() {
		case SubscriptionCreatedEvent:
			err = h.handleSubscriptionCreatedEvent(msg)
		default:
			h.logger.Infof("skipping message with subject %v ...", msg.Subject())

			return
		}

		if err != nil {
			h.logger.Error(err)

			err = msg.Nak()
			if err != nil {
				h.logger.Errorf("failed to negatively acknowledge message: %v", err)
			}

			return
		}

		err = msg.Ack()
		if err != nil {
			h.logger.Errorf("failed to acknowledge message: %v", err)
		}
	})
}

func (h *eventHandler) HandleCommands() error {
	return h.broker.ConsumeCommand(func(msg broker.ConsumedMessage) {
		var err error

		switch msg.Subject() {
		case SendDispatchCommand:
			err = h.handleSendDispatchCommand(msg)
		default:
			h.logger.Infof("skipping message with subject %v ...", msg.Subject())

			return
		}

		if err != nil {
			h.logger.Error(err)

			err = msg.Nak()
			if err != nil {
				h.logger.Errorf("failed to negatively acknowledge message: %v", err)
			}

			return
		}

		err = msg.Ack()
		if err != nil {
			h.logger.Errorf("failed to acknowledge message: %v", err)
		}
	})
}
