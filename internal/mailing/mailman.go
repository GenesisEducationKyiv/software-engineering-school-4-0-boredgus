package mailing

import "github.com/go-mail/mail"

type Email struct {
	From     string
	To       []string
	ReplyTo  string
	Subject  string
	HTMLBody string
}

type Mailman interface {
	Send(email Email) error
}

type mailman struct {
	dialer mail.Dialer
}

type SMTPParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewMailman(params SMTPParams) Mailman {
	return &mailman{
		dialer: *mail.NewDialer(params.Host, params.Port, params.Username, params.Password),
	}
}

func (m *mailman) Send(e Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", e.From)
	// TODO: find a way to send emails separately (without random subscriber knowing emails of other subscribers)
	msg.SetHeader("To", e.To...)
	msg.SetHeader("Subject", e.Subject)
	msg.SetBody("text/html", e.HTMLBody)
	return m.dialer.DialAndSend(msg)
}
