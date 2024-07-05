package notifier

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"

type (
	Email struct {
		To       []string
		Subject  string
		HTMLBody string
	}

	Mailman interface {
		Send(email Email) error
	}

	emailNotifier struct {
		notifier service.Notifier
		mailman  Mailman
	}
)

func NewEmailNotifier(notifier service.Notifier, mailman Mailman) *emailNotifier {
	return &emailNotifier{
		notifier: notifier,
		mailman:  mailman,
	}
}

func (n *emailNotifier) Notify(notification service.Notification) error {
	return nil
}
