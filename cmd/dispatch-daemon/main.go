package main

import (
	"bytes"
	"flag"
	"html/template"
	"os"
	"strconv"
	"subscription-api/config"
	"subscription-api/internal/mailing"
)

var (
	envFiles = flag.String("env", "dev.env", "list of env files separated with coma (e.g. '.env,prod.env')")
)

func init() {
	flag.Parse()
	config.InitEnvVariables(*envFiles)

}

func main() {
	config.InitLogger(config.DevMode)
	logger := config.InitLogger(config.Mode(os.Getenv("MODE")))

	from := os.Getenv("MAILMAL_EMAIL")
	port, err := strconv.Atoi(os.Getenv("MAILMAL_PORT"))
	if err != nil {
		logger.Errorf("invalid smtp server port: %v", err)

		return
	}
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
		Host:     os.Getenv("MAILMAL_HOST"),
		Port:     port,
		Username: from,
		Password: os.Getenv("MAILMAL_PASSWORD")}).
		Send(mailing.Email{
			From:     from,
			To:       []string{"daha@gmail.com"},
			ReplyTo:  from,
			Subject:  "Daily USD-UAH exchange rate",
			HTMLBody: buffer.String(),
		}))
}
