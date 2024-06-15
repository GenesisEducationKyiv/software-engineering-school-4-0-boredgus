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

func (m *mailman) Send(e Email) error {
	msgs := make([]*mail.Message, 0, len(e.To))

	for _, target := range e.To {
		msg := mail.NewMessage()
		msg.SetHeader("From", m.author)
		msg.SetHeader("To", target)
		msg.SetHeader("Subject", e.Subject)
		msg.SetBody("text/html", e.HTMLBody)
		msgs = append(msgs, msg)
	}

	return m.dialer.DialAndSend(msgs...)
}
