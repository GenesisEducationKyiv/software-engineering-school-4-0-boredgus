package app

import (
	"time"

	broker_msgs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
)

func ProtoToSubscription(msg *broker_msgs.SubscriptionMessage) *entities.Subscription {
	return &entities.Subscription{
		DispatchID:  msg.Payload.DispatchID,
		BaseCcy:     msg.Payload.BaseCcy,
		TargetCcies: msg.Payload.TargetCcies,
		Email:       msg.Payload.Email,
		SendAt:      msg.Payload.SendAt.AsTime()}
}

func ProtoToDispatchDetailsNotification(msg *broker_msgs.SubscriptionMessage, ntfctnType service.NotificationType) *service.Notification {
	return &service.Notification{
		Type: ntfctnType,
		Data: service.NotificationData{
			Emails: []string{msg.Payload.Email},
			Payload: service.SubscriptionData{
				BaseCcy:     msg.Payload.BaseCcy,
				TargetCcies: msg.Payload.TargetCcies,
				SendAt:      msg.Payload.SendAt.AsTime().UTC().Format(time.TimeOnly),
			},
		},
	}
}

func ProtoToCurrencyDispatchNotification(msg *broker_msgs.SendDispatchCommand) *service.Notification {
	return &service.Notification{
		Type: service.SendExchangeRates,
		Data: service.NotificationData{
			Emails: msg.Payload.Emails,
			Payload: service.CurrencyDispatchData{
				BaseCcy: msg.Payload.BaseCcy,
				Rates:   msg.Payload.Rates,
			},
		},
	}
}
