package entities

import "time"

type SubscriptionStatus int

func (s SubscriptionStatus) IsActive() bool {
	return s == SubscriptionStatusActive
}

const (
	SubscriptionStatusActive SubscriptionStatus = iota + 1
	SubscriptionStatusCancelled
)

type Subscription struct {
	DispatchID  string
	Email       string
	BaseCcy     string
	TargetCcies []string
	SendAt      time.Time
	Status      SubscriptionStatus
}
