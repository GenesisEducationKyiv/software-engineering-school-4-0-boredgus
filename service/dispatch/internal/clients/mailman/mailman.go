package mailman

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/go-mail/mail"
)

type Dialer interface {
	DialAndSend(m ...*mail.Message) error
}

type mailman struct {
	dialer Dialer
	author string
}

type SMTPParams struct {
	Host     string
	Port     int
	Email    string
	Name     string
	Password string
}

func NewMailman(params SMTPParams) *mailman {
	return &mailman{
		dialer: mail.NewDialer(params.Host, params.Port, params.Email, params.Password),
		author: fmt.Sprintf(`"%s" <%s>`, params.Name, params.Email),
	}
}

func (m *mailman) Send(email service.Email) error {
	msgs := make([]*mail.Message, 0, len(email.To))

	for _, target := range email.To {
		msg := mail.NewMessage()
		msg.SetHeader("From", m.author)
		msg.SetHeader("To", target)
		msg.SetHeader("Subject", email.Subject)
		msg.SetBody("text/html", email.HTMLBody)
		msgs = append(msgs, msg)
	}

	return m.dialer.DialAndSend(msgs...)
}
