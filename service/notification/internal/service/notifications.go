package service

type NotificationType int

const (
	SubscriptionCreated NotificationType = iota
	SendExchangeRates
	SubscriptionCancelled
)

type NotificationData struct {
	Emails  []string
	Payload interface{}
}

type Notification struct {
	Type NotificationType
	Data NotificationData
}

type SubscriptionData struct {
	BaseCcy     string
	TargetCcies []string
	SendAt      string
}

type CurrencyDispatchData struct {
	BaseCcy string
	Rates   map[string]float64
}
