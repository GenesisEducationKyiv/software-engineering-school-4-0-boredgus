package service

import (
	"context"

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
		GetStatusOfSubscription(ctx context.Context, args SubscriptionData) (entities.SubscriptionStatus, error)
		CreateSubscription(ctx context.Context, args SubscriptionData) error
		UpdateSubscriptionStatus(ctx context.Context, args SubscriptionData, status entities.SubscriptionStatus) error
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
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrUniqueViolation = errors.New("unique violation")
	ErrAlreadyExists   = errors.New("subscription already exists")
)

type dispatchService struct {
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
}

func NewDispatchService(
	userRepo UserRepo,
	subRepo SubRepo,
	dispatchRepo DispatchRepo,
) *dispatchService {
	return &dispatchService{
		userRepo:     userRepo,
		subRepo:      subRepo,
		dispatchRepo: dispatchRepo,
	}
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) (*entities.Subscription, error) {
	dispatchData, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return nil, err
	}

	subscription := dispatchData.ToSubscription(email, entities.SubscriptionStatusActive)
	subData := SubscriptionData{Email: email, DispatchID: dispatchId}

	status, err := s.subRepo.GetStatusOfSubscription(ctx, subData)
	if errors.Is(err, ErrNotFound) {
		err = s.createSubscription(ctx, subData)
		if err != nil {
			return nil, err
		}

		return subscription, nil
	}
	if err != nil {
		return nil, err
	}

	if status.IsActive() {
		return nil, ErrAlreadyExists
	}

	err = s.subRepo.UpdateSubscriptionStatus(ctx, subData, entities.SubscriptionStatusActive)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *dispatchService) createSubscription(ctx context.Context, subData SubscriptionData) error {
	if err := s.userRepo.CreateUser(ctx, subData.Email); err != nil && !errors.Is(err, ErrUniqueViolation) {
		return err
	}
	if err := s.subRepo.CreateSubscription(ctx, subData); err != nil {
		return err
	}

	return nil
}

func (s *dispatchService) UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) (*entities.Subscription, error) {
	dispatch, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return nil, err
	}

	err = s.subRepo.UpdateSubscriptionStatus(ctx,
		SubscriptionData{Email: email, DispatchID: dispatchId},
		entities.SubscriptionStatusCancelled,
	)
	if err != nil {
		return nil, err
	}

	return dispatch.ToSubscription(email, entities.SubscriptionStatusCancelled), nil
}
