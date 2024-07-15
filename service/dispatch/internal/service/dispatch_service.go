package service

import (
	"context"
	"time"

	"errors"

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
		Email, DispatchID string
	}
	SubRepo interface {
		CreateSubscription(ctx context.Context, args SubscriptionData) error
		UpdateSubscriptionStatus(ctx context.Context, args SubscriptionData, status SubscriptionStatus) error
	}

	DispatchRepo interface {
		GetDispatchByID(ctx context.Context, dispatchId string) (entities.CurrencyDispatch, error)
		GetSubscribersOfDispatch(ctx context.Context, dispatchId string) ([]string, error)
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
		CancelSubscription(sub Subscription)
	}
)

type SubscriptionStatus int

func (s SubscriptionStatus) IsActive() bool {
	return s == SubscriptionStatusActive
}

const (
	SubscriptionStatusActive SubscriptionStatus = iota + 1
	SubscriptionStatusCancelled
)

var (
	InvalidArgumentErr = errors.New("invalid argument")
	NotFoundErr        = errors.New("not found")
	UniqueViolationErr = errors.New("unique violation")
	AlreadyExistsErr   = errors.New("subscription already exists")
)

type dispatchService struct {
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
	broker       Broker
}

func NewDispatchService(
	userRepo UserRepo,
	subRepo SubRepo,
	dispatchRepo DispatchRepo,
	broker Broker,
) *dispatchService {
	return &dispatchService{
		userRepo:     userRepo,
		subRepo:      subRepo,
		dispatchRepo: dispatchRepo,
		broker:       broker,
	}
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchID string) error {
	dispatchData, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchID)
	if err != nil {
		return err
	}

	if err := s.userRepo.CreateUser(ctx, email); err != nil && !errors.Is(err, UniqueViolationErr) {
		return err
	}

	err = s.subRepo.CreateSubscription(ctx, SubscriptionData{Email: email, DispatchID: dispatchID})
	if errors.Is(err, UniqueViolationErr) {
		return AlreadyExistsErr
	}
	if err != nil {
		return err
	}

	s.broker.CreateSubscription(DispatchToSubscription(dispatchData, email))

	return nil
}

func (s *dispatchService) UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) error {
	dispatch, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	err = s.subRepo.UpdateSubscriptionStatus(ctx,
		SubscriptionData{Email: email, DispatchID: dispatchId},
		SubscriptionStatusCancelled,
	)
	if err != nil {
		return err
	}

	s.broker.CancelSubscription(DispatchToSubscription(dispatch, email))

	return nil
}

func DispatchToSubscription(dispatchData entities.CurrencyDispatch, email string) Subscription {
	return Subscription{
		DispatchID:  dispatchData.ID,
		Email:       email,
		BaseCcy:     dispatchData.Details.BaseCurrency,
		TargetCcies: dispatchData.Details.TargetCurrencies,
		SendAt:      dispatchData.SendAt,
	}
}
