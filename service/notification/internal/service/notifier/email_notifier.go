package notifier

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier/emails"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier/emails/templates"
)

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
	if len(notification.Data.Emails) > 0 {
		emailTemplate, err := templates.NotificationTypeToTemplate(notification.Type)
		if err != nil {
			return err
		}

		htmlBody, err := emails.ParseHTMLTemplate(emailTemplate.Name, notification.Data.Payload)
		if err != nil {
			return err
		}

		err = n.mailman.Send(Email{
			To:       notification.Data.Emails,
			Subject:  emailTemplate.Subject,
			HTMLBody: string(htmlBody),
		})
		if err != nil {
			return err
		}
	}

	return n.notifier.Notify(notification)
}
