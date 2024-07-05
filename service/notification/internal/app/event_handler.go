package app

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
)

type (
	Broker interface {
		Publish(subject string, data []byte) error
		Consume(subject, queue string, handler func(msg *broker.ConsumedMessage)) error
	}
	Scheduler interface {
	}
	NotificationService interface {
	}

	eventHandler struct {
		broker    Broker
		scheduler Scheduler
		service   NotificationService
	}
)

func NewEventHandler(broker Broker,
	scheduler Scheduler,
	service NotificationService,
) *eventHandler {
	return &eventHandler{
		scheduler: scheduler,
		broker:    broker,
		service:   service,
	}
}

func (h *eventHandler) CreateSubscription() error {
	return nil
}
func (h *eventHandler) DeleteSubscription() error {
	return nil
}
func (h *eventHandler) SendNotification() error {
	return nil
}
