package ds

import (
	"context"
	"errors"
	"subscription-api/internal/db"
	"subscription-api/internal/services"
)

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
