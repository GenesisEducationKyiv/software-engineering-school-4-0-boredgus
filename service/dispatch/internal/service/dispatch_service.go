package service

import (
	"context"
	"time"

	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
)

type (
	DispatchData struct {
		Id                 string
		Label              string
		SendAt             string
		CountOfSubscribers int
	}

	UserRepo interface {
		CreateUser(ctx context.Context, email string) error
	}

	SubscriptionData struct {
		Email, Dispatch string
	}
	SubRepo interface {
		CreateSubscription(ctx context.Context, args SubscriptionData) error
	}

	DispatchRepo interface {
		GetDispatchByID(ctx context.Context, dispatchId string) (entities.CurrencyDispatch, error)
		GetSubscribersOfDispatch(ctx context.Context, dispatchId string) ([]string, error)
		GetAllDispatches(ctx context.Context) ([]DispatchData, error)
	}

	Email struct {
		To       []string
		Subject  string
		HTMLBody string
	}

	Mailman interface {
		Send(email Email) error
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

var (
	InvalidArgumentErr = errors.New("invalid argument")
	NotFoundErr        = errors.New("not found")
	UniqueViolationErr = errors.New("unique violation")
)

type dispatchService struct {
	log          config.Logger
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
	broker       Broker
}

func NewDispatchService(
	logger config.Logger,
	userRepo UserRepo,
	subRepo SubRepo,
	dispatchRepo DispatchRepo,
	broker Broker,
) *dispatchService {
	return &dispatchService{
		userRepo:     userRepo,
		subRepo:      subRepo,
		dispatchRepo: dispatchRepo,
		log:          logger,
		broker:       broker,
	}
}

func (s *dispatchService) GetAllDispatches(ctx context.Context) ([]DispatchData, error) {
	return s.dispatchRepo.GetAllDispatches(ctx)
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	dispatchData, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	if err = s.userRepo.CreateUser(ctx, email); err != nil && !errors.Is(err, UniqueViolationErr) {
		return err
	}

	if err = s.subRepo.CreateSubscription(ctx, SubscriptionData{Email: email, Dispatch: dispatchId}); err != nil {
		return err
	}

	s.broker.CreateSubscription(Subscription{
		DispatchID:  dispatchId,
		Email:       email,
		BaseCcy:     dispatchData.Details.BaseCurrency,
		TargetCcies: dispatchData.Details.TargetCurrencies,
		SendAt:      dispatchData.SendAt,
	})

	return nil
}
