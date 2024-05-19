package ds

import (
	"context"
	"fmt"
	"strings"
	"subscription-api/internal/db"
	"subscription-api/internal/services"
)

const getDispatchQ string = `
	select cd.u_id, cd.base_currency , cd.target_currencies, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	where cd.u_id = $1
	group by cd.id;
`

func (s *currencyDispatchStore) GetDispatch(ctx context.Context, dispatchId string) (DispatchData, error) {
	var d DispatchData
	row := s.store.DB().QueryRow(getDispatchQ, dispatchId)
	err := row.Err()
	if s.IsError(err, db.InvalidTextRepresentation) {
		return d, fmt.Errorf("%w: incorrect format for uuid", services.InvalidArgumentErr)
	}
	var targetCurrencies string
	if err := row.Scan(&d.Id, &d.Details.BaseCurrency, &targetCurrencies, &d.SendAt, &d.CountOfSubscribers); err != nil {
		return d, fmt.Errorf("%w: dispatch with such id does not exists", services.NotFoundErr)
	}
	d.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")
	return d, nil
}
