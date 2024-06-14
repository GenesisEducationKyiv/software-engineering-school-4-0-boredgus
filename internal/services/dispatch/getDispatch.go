package ds

import (
	"context"
	db "subscription-api/internal/db/dispatch"
	"subscription-api/internal/entities"
)

func (s *dispatchService) GetDispatch(ctx context.Context, dispatchId string) (DispatchInfo, error) {
	var dispatch DispatchInfo
	err := s.store.WithTx(ctx, func(dq db.DispatchQueries) error {
		d, err := dq.GetDispatch(ctx, dispatchId)
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
