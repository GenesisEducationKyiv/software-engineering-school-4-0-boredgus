package ds

import (
	"context"
	db "subscription-api/internal/db"
	"subscription-api/internal/entities"
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
