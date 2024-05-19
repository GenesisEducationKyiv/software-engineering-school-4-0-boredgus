package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
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

	from := "daha.kyiv@gmail.com"
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
	fmt.Println(mailing.NewMailman(mailing.SMTPParams{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: from,
		Password: "guze dokh umzh ulvs"}).
		Send(mailing.Email{
			From:     from,
			To:       []string{"kefirchi@ukr.net"},
			ReplyTo:  from,
			Subject:  "Daily USD-UAH exchange rate",
			HTMLBody: buffer.String(),
		}))
}
