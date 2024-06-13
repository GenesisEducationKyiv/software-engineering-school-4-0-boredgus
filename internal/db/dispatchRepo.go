package db

import (
	"context"
	"fmt"
	"strings"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type DispatchRepo struct{}

func NewDispatchRepo() *DispatchRepo {
	return &DispatchRepo{}
}

const getDispatchQ string = `
	select cd.u_id, cd.base_currency , cd.target_currencies, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	where cd.u_id = $1
	group by cd.id;
`

type DispatchData struct {
	entities.Dispatch[entities.CurrencyDispatchDetails]
	CountOfSubscribers int
}

func (s *DispatchRepo) GetByID(ctx context.Context, d DB, dispatchId string) (DispatchData, error) {
	var data DispatchData
	row := d.DB().QueryRow(getDispatchQ, dispatchId)
	err := row.Err()
	if d.IsError(err, InvalidTextRepresentation) {
		return data, fmt.Errorf("%w: incorrect format for uuid", services.InvalidArgumentErr)
	}
	var targetCurrencies string
	if err := row.Scan(&data.Id, &data.Details.BaseCurrency, &targetCurrencies, &data.SendAt, &data.CountOfSubscribers); err != nil {
		return data, fmt.Errorf("%w: dispatch with such id does not exists", services.NotFoundErr)
	}
	data.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")

	return data, nil
}
