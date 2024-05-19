package ds

import (
	"context"
	ds "subscription-api/internal/db/dispatch"
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

type dispatchService struct {
	store ds.DispatchStore
}

func NewDispatchService(s ds.DispatchStore) DispatchService {
	return &dispatchService{store: s}
}
