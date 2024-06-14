package ds

import (
	"context"
	"errors"
	db "subscription-api/internal/db"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type DispatchInfo struct {
	entities.Dispatch[entities.CurrencyDispatchDetails]
	CountOfSubscribers int
}

type DispatchService interface {
	SubscribeForDispatch(ctx context.Context, email, dispatch string) error
	SendDispatch(ctx context.Context, dispatch string) error
	GetDispatch(ctx context.Context, dispatch_id string) (DispatchInfo, error)
}
type UserRepo interface {
	CreateUser(ctx context.Context, db db.DB, email string) error
}

type SubRepo interface {
	CreateSubscription(ctx context.Context, db db.DB, args db.SubscriptionData) error
}

type DispatchRepo interface {
	GetByID(ctx context.Context, db db.DB, dispatchId string) (db.DispatchData, error)
}

type dispatchService struct {
	store        db.Store
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
}

func NewDispatchService(s db.Store) DispatchService {
	return &dispatchService{
		store:        s,
		userRepo:     db.NewUserRepo(),
		subRepo:      db.NewSubRepo(),
		dispatchRepo: db.NewDispatchRepo(),
	}
}

func (s *dispatchService) GetDispatch(ctx context.Context, dispatchId string) (DispatchInfo, error) {
	var dispatch DispatchInfo
	err := s.store.WithTx(ctx, func(db db.DB) error {
		d, err := s.dispatchRepo.GetByID(ctx, db, dispatchId)
		if err == nil {
			dispatch = DispatchInfo{
				Dispatch: entities.Dispatch[entities.CurrencyDispatchDetails]{
					Id:     d.Id,
					SendAt: d.SendAt,
					Details: entities.CurrencyDispatchDetails{
						BaseCurrency:     d.Details.BaseCurrency,
						TargetCurrencies: d.Details.TargetCurrencies,
					},
				},
				CountOfSubscribers: d.CountOfSubscribers,
			}
		}

		return err
	})

	return dispatch, err
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	return s.store.WithTx(ctx, func(d db.DB) error {
		_, err := s.dispatchRepo.GetByID(ctx, d, dispatchId)
		if err != nil {
			return err
		}

		if err = s.userRepo.CreateUser(ctx, d, email); err != nil && !errors.Is(err, services.UniqueViolationErr) {
			return err
		}

		return s.subRepo.CreateSubscription(ctx, d, db.SubscriptionData{Email: email, Dispatch: dispatchId})
	})
}

func (s *dispatchService) SendDispatch(ctx context.Context, dispatch string) error {
	return nil
}
