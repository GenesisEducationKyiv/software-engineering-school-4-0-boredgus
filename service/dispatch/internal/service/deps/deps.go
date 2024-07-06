package deps

import (
	"context"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
)

type (
	DispatchData struct {
		Id                 string
		Label              string
		SendAt             string
		CountOfSubscribers int
	}
	SubscriptionData struct {
		Email, Dispatch string
	}

	UserRepo interface {
		CreateUser(ctx context.Context, email string) error
	}

	SubRepo interface {
		CreateSubscription(ctx context.Context, args SubscriptionData) error
	}

	DispatchRepo interface {
		GetDispatchByID(ctx context.Context, dispatchId string) (entities.CurrencyDispatch, error)
		GetSubscribersOfDispatch(ctx context.Context, dispatchId string) ([]string, error)
		GetAllDispatches(ctx context.Context) ([]DispatchData, error)
	}

	Subscription struct {
		DispatchID  string
		Email       string
		BaseCcy     string
		TargetCcies []string
		SendAt      time.Time
	}

	Broker interface {
		CreateSubscription(sub Subscription)
	}
)
