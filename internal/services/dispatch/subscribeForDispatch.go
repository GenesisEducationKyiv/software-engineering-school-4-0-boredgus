package ds

import (
	"context"
	"errors"
	db "subscription-api/internal/db/dispatch"
	"subscription-api/internal/services"
)

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	return s.store.WithTx(ctx, func(dq db.DispatchQueries) error {
		_, err := dq.GetDispatch(ctx, dispatchId)
		if err != nil {
			return err
		}
		err = dq.CreateSubscriber(ctx, email)
		if err != nil && !errors.Is(err, services.UniqueViolationErr) {
			return err
		}

		return dq.SubscribeFor(ctx, db.SubscribeForParams{Email: email, Dispatch: dispatchId})
	})
}
