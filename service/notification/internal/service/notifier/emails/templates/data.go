package templates

import "fmt"

type Template string

const (
	SubscriptionCreated  Template = "subscription_created"
	ExchangeRateTemplate Template = "exchange_rate"
)

var templateToSubjectMapper = map[Template]string{
	SubscriptionCreated:  "You subscribed to exchange rate dispatch!",
	ExchangeRateTemplate: "Exchange rate",
}

func TemplateToSubject(templateName string) (string, error) {
	template := Template(templateName)
	subjectTemplate, ok := templateToSubjectMapper[template]
	if !ok {
		return "", fmt.Errorf("provided unsupported email template name '%s'", template)
	}

	return subjectTemplate, nil
}
