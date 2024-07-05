package service

type NotificationType int

const (
	SubscriptionCreated NotificationType = iota
	SendExchangeRates
)

type NotificationData struct {
	Emails  []string
	Payload interface{}
}

type Notification struct {
	Type NotificationType
	Data NotificationData
}

type SubscriptionDetails struct {
	BaseCcy     string
	TargetCcies []string
	SendAt      string
}

type CurrencyDispatch struct {
	BaseCcy string
	Rates   map[string]float64
}
