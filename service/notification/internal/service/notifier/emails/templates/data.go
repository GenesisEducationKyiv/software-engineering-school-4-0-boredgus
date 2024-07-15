package templates

import (
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
)

type Template struct {
	Name, Subject string
}

var notificationTypeToTemplate = map[service.NotificationType]Template{
	service.SubscriptionCreated: {
		Name:    "subscription_created",
		Subject: "You subscribed to exchange rate dispatch!",
	},
	service.SendExchangeRates: {
		Name:    "exchange_rate",
		Subject: "Exchange rates",
	},
	service.SubscriptionCancelled: {
		Name:    "subscription_cancelled",
		Subject: "Subscription cancellation",
	},
}

var UnknownNotificationTypeErr = errors.New("unknown notification type")

func NotificationTypeToTemplate(ntfctionType service.NotificationType) (Template, error) {
	subjectTemplate, ok := notificationTypeToTemplate[ntfctionType]
	if !ok {
		return Template{}, fmt.Errorf("%w: email template for '%v' notification is not configured", UnknownNotificationTypeErr, ntfctionType)
	}

	return subjectTemplate, nil
}
