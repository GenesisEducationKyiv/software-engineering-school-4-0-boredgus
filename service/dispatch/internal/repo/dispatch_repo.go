package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
)

type dispatchRepo struct {
	db DB
}

func NewDispatchRepo(db DB) *dispatchRepo {
	return &dispatchRepo{db: db}
}

const getDispatchByIdQ string = `
	select cd.u_id, cd.label, cd.template_name, cd.base_currency , cd.target_currencies, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	where cd.u_id = $1
	group by cd.id;
`

func (r *dispatchRepo) GetDispatchByID(ctx context.Context, dispatchId string) (entities.CurrencyDispatch, error) {
	var dispatch entities.CurrencyDispatch
	row := r.db.QueryRowContext(ctx, getDispatchByIdQ, dispatchId)
	err := row.Err()
	if err != nil && r.db.IsError(err, InvalidTextRepresentation) {
		return dispatch, fmt.Errorf("%w: incorrect format for uuid", service.ErrInvalidArgument)
	}
	if err != nil {
		return dispatch, err
	}
	var targetCurrencies string
	if err := row.Scan(&dispatch.ID, &dispatch.Label, &dispatch.TemplateName, &dispatch.Details.BaseCurrency, &targetCurrencies, &dispatch.SendAt, &dispatch.CountOfSubscribers); err != nil {
		return dispatch, fmt.Errorf("%w: dispatch with such id does not exists", service.ErrNotFound)
	}
	dispatch.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")

	return dispatch, nil
}

const getSubscribersOfDispatchQ string = `
	select u.email
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	left join subs."users" u 
	on cs.user_id = u.id
	where cd.u_id = $1 and u.email is not null and cs.status != (
		select id
		from subs."subscription_status"
		where label = 'cancelled'
	);
`

func (r *dispatchRepo) GetSubscribersOfDispatch(ctx context.Context, dispatchId string) ([]string, error) {
	var result []string
	rows, err := r.db.QueryContext(ctx, getSubscribersOfDispatchQ, dispatchId)
	if r.db.IsError(err, InvalidTextRepresentation) {
		return result, fmt.Errorf("%w: incorrect format for uuid", service.ErrInvalidArgument)
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return result, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, email)
	}

	return result, nil
}
