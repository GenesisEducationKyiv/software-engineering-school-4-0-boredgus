package ds

import (
	"context"
	"fmt"
	"subscription-api/internal/db"
	"subscription-api/internal/services"
)

const subscribeForQ string = `
	insert into subs."currency_subscriptions" (user_id, dispatch_id)
	select *
	from 
		(select u.id user_id
		from subs."users" as u
		where u.email = $1) u,
		(select cd.id dispatch_id
		from subs.currency_dispatches cd
		where cd.u_id = $2) d
	;
`

func (s *currencyDispatchStore) SubscribeFor(ctx context.Context, params SubscribeForParams) error {
	_, err := s.DB().ExecContext(ctx, subscribeForQ, params.Email, params.Dispatch)
	if s.IsError(err, db.UniqueViolation) {
		return fmt.Errorf("%w: user has already subscribed for this dispatch", services.UniqueViolationErr)
	}

	return err
}
