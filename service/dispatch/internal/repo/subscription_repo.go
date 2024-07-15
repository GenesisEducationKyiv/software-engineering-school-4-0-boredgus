package repo

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
)

type subscriptionRepo struct {
	db DB
}

func NewSubRepo(db DB) *subscriptionRepo {
	return &subscriptionRepo{db: db}
}

const createSubscriptionQ string = `
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

func (r *subscriptionRepo) CreateSubscription(ctx context.Context, args service.SubscriptionData) error {
	_, err := r.db.ExecContext(ctx, createSubscriptionQ, args.Email, args.DispatchID)
	if r.db.IsError(err, UniqueViolation) {
		return fmt.Errorf("%w: user has already subscribed for this dispatch", service.UniqueViolationErr)
	}

	return err
}

const updateSubscriptionStatusQ = `
	update subs."currency_subscriptions"
	set status = $1
	from 
		subs."users" users,
		subs."currency_dispatches" cds
	where 
		users.email = $2 and user_id = users.id and
		cds.u_id = $3 and dispatch_id = cds.id;
`

func (r *subscriptionRepo) UpdateSubscriptionStatus(ctx context.Context, sub service.SubscriptionData, status service.SubscriptionStatus) error {
	rows, err := r.db.ExecContext(ctx, updateSubscriptionStatusQ, status, sub.Email, sub.DispatchID)
	if err != nil {
		return err
	}

	count, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return service.NotFoundErr
	}

	return err
}
