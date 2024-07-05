package service

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
)

type (
	Notifier interface {
		Notify(Notification) error
	}
	CurrencyService interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}
	notificationService struct {
		logger          config.Logger
		notifier        Notifier
		currencyService CurrencyService
	}
)

func NewNotificationService(
	logger config.Logger,
	notifier Notifier,
	currencyService CurrencyService,
) *notificationService {
	return &notificationService{
		logger:          logger,
		notifier:        notifier,
		currencyService: currencyService,
	}
}

func (s *notificationService) CreateSubscription(
	ctx context.Context,
	notification CreateSubscriptionNotification,
) error {

	return nil
}

func (s *notificationService) SendExchangeRate(
	ctx context.Context,
	notification ExchangeRateNotification,
) error {

	return nil
}
