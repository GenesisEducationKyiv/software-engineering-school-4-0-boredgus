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
		GetStatusOfSubscription(ctx context.Context, args SubscriptionData) (SubscriptionStatus, error)
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
		Status      SubscriptionStatus
	}

	Broker interface {
		Publish(msg interface{})
	}
)

type SubscriptionStatus int64

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

	subData := SubscriptionData{Email: email, DispatchID: dispatchID}
	status, err := s.subRepo.GetStatusOfSubscription(ctx, subData)
	if err != nil && errors.Is(err, NotFoundErr) {
		return s.createSubscription(ctx, subData, dispatchData)
	}
	if err != nil {
		return err
	}

	if status.IsActive() {
		return AlreadyExistsErr
	}

	err = s.subRepo.UpdateSubscriptionStatus(ctx, subData, SubscriptionStatusActive)
	if err != nil {
		return err
	}
	s.broker.Publish(DispatchToSubscription(dispatchData, email, SubscriptionStatusActive))

	return nil
}

func (s *dispatchService) createSubscription(ctx context.Context, subData SubscriptionData, dispatch entities.CurrencyDispatch) error {
	if err := s.userRepo.CreateUser(ctx, subData.Email); err != nil && !errors.Is(err, UniqueViolationErr) {
		return err
	}
	if err := s.subRepo.CreateSubscription(ctx, subData); err != nil {
		return err
	}

	s.broker.Publish(DispatchToSubscription(dispatch, subData.Email, SubscriptionStatusActive))

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

	s.broker.Publish(DispatchToSubscription(dispatch, email, SubscriptionStatusCancelled))

	return nil
}

func DispatchToSubscription(
	dispatchData entities.CurrencyDispatch,
	email string,
	status SubscriptionStatus,
) Subscription {
	return Subscription{
		DispatchID:  dispatchData.ID,
		Email:       email,
		BaseCcy:     dispatchData.Details.BaseCurrency,
		TargetCcies: dispatchData.Details.TargetCurrencies,
		SendAt:      dispatchData.SendAt,
		Status:      status,
	}
}
