package app

// import (
// 	"context"

// 	nats_msgs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/notifier/internal/broker/gen"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/notifier/internal/clients/smtp"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/notifier/internal/clients/smtp/emails"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/notifier/internal/clients/smtp/emails/templates"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/notifier/internal/config"
// )

// type (
// 	Mailman interface {
// 		Send(email smtp.Email) error
// 	}

// 	Message interface {
// 		Data() (*nats_msgs.SendNotificationMessage, error)
// 		Ack(context.Context) error
// 		Nak(context.Context) error
// 	}

// 	Consumer interface {
// 		Consume(subject, queue string, handler func(msg Message)) error
// 	}
// 	app struct {
// 		mailman Mailman
// 		broker  Consumer
// 		logger  config.Logger
// 	}
// )

// func NewApp(mailman Mailman, broker Consumer, logger config.Logger) *app {
// 	return &app{
// 		mailman: mailman,
// 		broker:  broker,
// 		logger:  logger,
// 	}
// }

// func (a *app) sendEmails(data *nats_msgs.Data) error {
// 	subject, err := templates.TemplateToSubject(data.TemplateName)
// 	if err != nil {
// 		return err
// 	}

// 	htmlContent, err := emails.ParseHTMLTemplate(data.TemplateName, data.Payload)
// 	if err != nil {
// 		return err
// 	}

// 	err = a.mailman.Send(smtp.Email{
// 		To:       data.Emails,
// 		Subject:  subject,
// 		HTMLBody: string(htmlContent),
// 	})
// 	if err != nil {
// 		a.logger.Errorf("failed to send email via smtp: %w", err)
// 	}

// 	return err
// }

// func (a *app) Run(subject string) {
// 	if err := a.broker.Consume(subject, "notifications", func(msg Message) {
// 		msgData, err := msg.Data()
// 		if err != nil {
// 			a.logger.Error(err)

// 			return
// 		}

// 		if len(msgData.Data.Emails) == 0 {
// 			return
// 		}

// 		ctx := context.Background()
// 		if err = a.sendEmails(msgData.Data); err != nil {
// 			err = msg.Nak(ctx)

// 			if err != nil {
// 				a.logger.Errorf("failed to negatively acknowledge message: %w", err)
// 			}

// 			return
// 		}

// 		if err := msg.Ack(ctx); err != nil {
// 			a.logger.Errorf("failed to acknowledge broker: %w", err)
// 		}
// 	}); err != nil {
// 		a.logger.Errorf("failed to consume queue messages: %w", err)
// 	}
// }
