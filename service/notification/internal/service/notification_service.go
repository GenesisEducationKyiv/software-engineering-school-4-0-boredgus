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

func (s *notificationService) SendSubscriptionDetails(ctx context.Context, notification Notification) error {
	err := s.notifier.Notify(notification)
	if err != nil {
		s.logger.Errorf("failed to send subscription details: %v", err)

		return err
	}

	s.logger.Infof("succefuly sent subscription details: %+v", notification.Data)

	return nil
}

func (s *notificationService) SendExchangeRates(ctx context.Context, notification Notification) error {
	err := s.notifier.Notify(notification)
	if err != nil {
		s.logger.Errorf("failed to send exchange rates: %v", err)

		return err
	}

	s.logger.Infof("succefuly sent exchange rates: %+v", notification.Data)

	return nil
}
