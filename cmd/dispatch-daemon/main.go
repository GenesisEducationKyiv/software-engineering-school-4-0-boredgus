package main

import (
	"bytes"
	"flag"
	"html/template"
	"subscription-api/cmd/dispatch-daemon/internal"
	"subscription-api/config"
	"subscription-api/internal/mailing"
)

var (
	envFiles = flag.String("env", "dev.env", "list of env files separated with coma (e.g. '.env,prod.env')")
)

func init() {
	flag.Parse()
	config.LoadEnvFile(*envFiles)

}

/*
This module is not implemented yet.
The logic written below is redundunt and not worth to pay attention to.
I'll implement it later.
*/
func main() {
	env := internal.Env()
	logger := config.InitLogger(env.Mode)

	data := struct {
		BaseCurrency   string
		TargetCurrency string
		ExchangeRate   float64
	}{
		BaseCurrency:   "USD",
		TargetCurrency: "UAH",
		ExchangeRate:   30.1232211,
	}
	var buffer bytes.Buffer
	if err := template.
		Must(template.ParseFiles("internal/mailing/emails/exchange_rate.html")).
		Execute(&buffer, data); err != nil {
		config.Log().Fatal("failed to execute template: ", err.Error())
	}
	logger.Info(mailing.NewMailman(mailing.SMTPParams{
		Host:     env.MailmanHost,
		Port:     env.MailmanPort,
		Username: env.MailmanEmail,
		Password: env.MailmanPassword}).
		Send(mailing.Email{
			From:     env.MailmanEmail,
			To:       []string{"daha@gmail.com"},
			ReplyTo:  env.MailmanEmail,
			Subject:  "Daily USD-UAH exchange rate",
			HTMLBody: buffer.String(),
		}))
}
