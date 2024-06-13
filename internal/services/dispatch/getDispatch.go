package ds

import (
	"context"
	"subscription-api/internal/db"
	"subscription-api/internal/entities"
)

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
