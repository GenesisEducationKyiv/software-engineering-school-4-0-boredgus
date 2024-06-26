package db

import (
	"context"
	"fmt"
	"subscription-api/internal/services"
)

type subRepo struct {
	db DB
}

func NewSubRepo(db DB) *subRepo {
	return &subRepo{db: db}
}

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

type SubscriptionData struct {
	Email, Dispatch string
}

func (r *subRepo) CreateSubscription(ctx context.Context, args SubscriptionData) error {
	_, err := r.db.ExecContext(ctx, subscribeForQ, args.Email, args.Dispatch)
	if r.db.IsError(err, UniqueViolation) {
		return fmt.Errorf("%w: user has already subscribed for this dispatch", services.UniqueViolationErr)
	}

	return err
}
