package testdata

import (
	"time"

	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var SubscriptionCreatedMessage = &messages.SubscriptionMessage{
	EventType: messages.EventType_SUBSCRIPTION_CREATED,
	Timestamp: timestamppb.New(time.Now().UTC()),
	Payload: &messages.Subscription{
		DispatchID:  "dispatch-id",
		Email:       "created_sub_email@gmail.com",
		BaseCcy:     "UAH",
		TargetCcies: []string{"USD", "EUR"},
		Status:      messages.SubscriptionStatus_CREATED,
		SendAt:      timestamppb.New(time.Now().UTC()),
	},
}

var SubscriptionCancelledMessage = &messages.SubscriptionMessage{
	EventType: messages.EventType_SUBSCRIPTION_CANCELLED,
	Timestamp: timestamppb.New(time.Now().UTC()),
	Payload: &messages.Subscription{
		DispatchID:  "dispatch-id",
		Email:       "cancelled_sub_email@gmail.com",
		BaseCcy:     "UAH",
		TargetCcies: []string{"USD", "EUR"},
		Status:      messages.SubscriptionStatus_CANCELLED,
		SendAt:      timestamppb.New(time.Now().UTC()),
	},
}
