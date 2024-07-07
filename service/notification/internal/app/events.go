package app

import (
	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
)

const (
	SubscriptionCreatedEvent   string = "events.subscription.created"
	SubscriptionCancelledEvent string = "events.subscription.cancelled"
	SubscriptionRenewedEvent   string = "events.subscription.renewed"
	SendDispatchCommand        string = "commands.send.dispatch"
)

var messageToNotificationMapper = map[messages.EventType]service.NotificationType{
	messages.EventType_SUBSCRIPTION_CREATED:   service.SubscriptionCreated,
	messages.EventType_SUBSCRIPTION_CANCELLED: service.SubscriptionCancelled,
	messages.EventType_SUBSCRIPTION_RENEWED:   service.SubscriptionRenewed,
}

// MessageTypeToNotificationType transforms message type to notification type.
func MessageTypeToNotificationType(t messages.EventType) service.NotificationType {
	return messageToNotificationMapper[t]
}
