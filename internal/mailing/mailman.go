package mailing

import (
	"fmt"

	"github.com/go-mail/mail"
)

type Email struct {
	To       []string
	Subject  string
	HTMLBody string
}

type Mailman interface {
	Send(e Email) error
}

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

func (m *mailman) Send(email Email) error {
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
