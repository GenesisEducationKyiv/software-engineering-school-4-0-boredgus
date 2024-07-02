package repo

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	service_errors "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"
)

type subscriptionRepo struct {
	db DB
}

func NewSubRepo(db DB) *subscriptionRepo {
	return &subscriptionRepo{db: db}
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

func (r *subscriptionRepo) CreateSubscription(ctx context.Context, args deps.SubscriptionData) error {
	_, err := r.db.ExecContext(ctx, subscribeForQ, args.Email, args.Dispatch)
	if r.db.IsError(err, UniqueViolation) {
		return fmt.Errorf("%w: user has already subscribed for this dispatch", service_errors.UniqueViolationErr)
	}

	return err
}
