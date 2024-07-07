package service

import (
	"context"
	"fmt"
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
		GetStatusOfSubscription(ctx context.Context, args SubscriptionData) (SubscriptionStatus, error)
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
		RenewSubscription(sub Subscription)
	}
)

type SubscriptionStatus int

func (s SubscriptionStatus) IsActive() bool {
	return s == SubscriptionCreatedStatus || s == SubscriptionRenewedStatus
}

const (
	SubscriptionCreatedStatus SubscriptionStatus = iota + 1
	SubscriptionCancelledStatus
	SubscriptionRenewedStatus
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

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	dispatchData, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	subscription := DispatchToSubscription(dispatchData, email)
	status, err := s.subRepo.GetStatusOfSubscription(ctx, SubscriptionData{Email: email, DispatchID: dispatchId})
	if errors.Is(err, NotFoundErr) {
		return s.createSubscription(ctx, subscription)
	}
	if err != nil {
		return err
	}
	if status.IsActive() {
		return AlreadyExistsErr
	}

	return s.renewSubscription(ctx, subscription)
}

func (s *dispatchService) createSubscription(ctx context.Context, sub Subscription) error {
	if err := s.userRepo.CreateUser(ctx, sub.Email); err != nil && !errors.Is(err, UniqueViolationErr) {
		return err
	}

	err := s.subRepo.CreateSubscription(ctx, SubscriptionData{Email: sub.Email, DispatchID: sub.DispatchID})
	if err != nil {
		return err
	}

	s.broker.CreateSubscription(sub)

	return nil
}

func (s *dispatchService) renewSubscription(ctx context.Context, sub Subscription) error {
	if err := s.subRepo.UpdateSubscriptionStatus(ctx,
		SubscriptionData{Email: sub.Email, DispatchID: sub.DispatchID},
		SubscriptionRenewedStatus,
	); err != nil {
		return err
	}

	s.broker.RenewSubscription(sub)

	return nil
}

func (s *dispatchService) UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) error {
	dispatch, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	subData := SubscriptionData{Email: email, DispatchID: dispatchId}
	status, err := s.subRepo.GetStatusOfSubscription(ctx, subData)
	if errors.Is(err, NotFoundErr) {
		return fmt.Errorf("%w: user is not subscired to provided dispatch", err)
	}
	if err != nil {
		return err
	}

	if status == SubscriptionCancelledStatus {
		return nil
	}

	err = s.subRepo.UpdateSubscriptionStatus(ctx, subData, SubscriptionCancelledStatus)
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
