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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Broker interface {
		PublishAsync(subject string, payload []byte) error
		ConsumeMessage(handler func(msg broker.ConsumedMessage)) error
	}

	NotificationService interface {
		SendSubscriptionDetails(ctx context.Context, notification service.Notification) error
		SendExchangeRates(ctx context.Context, notification service.Notification) error
	}

	DispatchStore interface {
		AddSubscription(ctx context.Context, sub entities.Subscription) error
	}

	Scheduler interface {
		AddSubscription(entities.Subscription, func(*entities.Dispatch))
	}

	Converter interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}

	eventHandler struct {
		broker        Broker
		converter     Converter
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
	RedeliveryDelay     time.Duration = 1 * time.Minute
)

func NewEventHandler(
	broker Broker,
	ccyConverter Converter,
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

func (h *eventHandler) invokeSendingOfDispatch(d *entities.Dispatch) {
	msg := broker_msgs.SendDispatchCommand{
		EventType: broker_msgs.EventType_SEND_DISPATCH,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload: &broker_msgs.Dispatch{
			Emails:  d.Emails,
			BaseCcy: d.BaseCcy,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	rates, err := h.converter.Convert(ctx, d.BaseCcy, d.TargetCcies)
	if err != nil {
		h.logger.Errorf("failed to get rates: %v", err)

		return
	}
	msg.Payload.Rates = rates

	marshalled, err := proto.Marshal(&msg)
	if err != nil {
		h.logger.Errorf("failed to marshal SendDispatchCommand: %v", err)

		return
	}

	if err := h.broker.PublishAsync(SendDispatchCommand, marshalled); err != nil {
		h.logger.Errorf("failed to emit SendDispatch commands: %v", err)
	}
}

func (h *eventHandler) handleSubscriptionCreatedEvent(msg broker.ConsumedMessage) error {
	var parsedMsg broker_msgs.SubscriptionMessage
	if err := proto.Unmarshal(msg.Data(), &parsedMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message from %s: %w", SubscriptionCreatedEvent, err)
	}

	sub := entities.Subscription{
		DispatchID:  parsedMsg.Payload.DispatchID,
		BaseCcy:     parsedMsg.Payload.BaseCcy,
		TargetCcies: parsedMsg.Payload.TargetCcies,
		Email:       parsedMsg.Payload.Email,
		SendAt:      parsedMsg.Payload.SendAt.AsTime(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	h.scheduler.AddSubscription(sub, h.invokeSendingOfDispatch)
	if err := h.dispatchStore.AddSubscription(ctx, sub); err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}

	if err := h.service.SendSubscriptionDetails(ctx, service.Notification{
		Type: service.SubscriptionCreated,
		Data: service.NotificationData{
			Emails: []string{parsedMsg.Payload.Email},
			Payload: service.SubscriptionData{
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
			Emails: parsedMsg.Payload.Emails,
			Payload: service.CurrencyDispatchData{
				BaseCcy: parsedMsg.Payload.BaseCcy,
				Rates:   parsedMsg.Payload.Rates,
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send subscription details: %w", err)
	}

	return nil
}

func (h *eventHandler) HandleMessages() error {
	return h.broker.ConsumeMessage(func(msg broker.ConsumedMessage) {
		var err error

		switch msg.Subject() {
		case SubscriptionCreatedEvent:
			err = h.handleSubscriptionCreatedEvent(msg)
		case SendDispatchCommand:
			err = h.handleSendDispatchCommand(msg)
		default:
			h.logger.Infof("skipping message with subject %v ...", msg.Subject())

			return
		}

		h.logger.Infof("handling message with subject %v ...", msg.Subject())

		if err != nil {
			h.logger.Errorf("failed to handle message: %v", err)

			err = msg.NakWithDelay(RedeliveryDelay)
			if err != nil {
				h.logger.Errorf("failed to negatively acknowledge message: %v", err)
			}

			return
		}

		err = msg.Ack()
		if err != nil {
			h.logger.Errorf("failed to acknowledge message: %v", err)
		}

		h.logger.Info("successfully handled message")
	})
}
