package notifier

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"

type baseNotifier struct{}

func NewBaseNotifier() *baseNotifier {
	return &baseNotifier{}
}

func (n *baseNotifier) Notify(notification service.Notification) error {
	return nil
}
