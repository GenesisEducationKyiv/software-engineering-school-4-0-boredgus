package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
)

type DispatchRepo struct{}

func NewDispatchRepo() *DispatchRepo {
	return &DispatchRepo{}
}

const getDispatchByIdQ string = `
	select cd.u_id, cd.label, cd.template_name, cd.base_currency , cd.target_currencies, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	where cd.u_id = $1
	group by cd.id;
`

func (r *DispatchRepo) GetDispatchByID(ctx context.Context, db DB, dispatchId string) (entities.CurrencyDispatch, error) {
	var dispatch entities.CurrencyDispatch
	row := db.DB().QueryRowContext(ctx, getDispatchByIdQ, dispatchId)
	err := row.Err()
	if err != nil && db.IsError(err, InvalidTextRepresentation) {
		return dispatch, fmt.Errorf("%w: incorrect format for uuid", services.InvalidArgumentErr)
	}
	if err != nil {
		return dispatch, err
	}
	var targetCurrencies string
	if err := row.Scan(&dispatch.Id, &dispatch.Label, &dispatch.TemplateName, &dispatch.Details.BaseCurrency, &targetCurrencies, &dispatch.SendAt, &dispatch.CountOfSubscribers); err != nil {
		return dispatch, fmt.Errorf("%w: dispatch with such id does not exists", services.NotFoundErr)
	}
	dispatch.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")

	return dispatch, nil
}

const getSubscribersOfDispatchQ string = `
	select u.email
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	left join users u 
	on cs.user_id = u.id
	where cd.u_id = $1;
`

func (r *DispatchRepo) GetSubscribersOfDispatch(ctx context.Context, d DB, dispatchId string) ([]string, error) {
	var result []string
	rows, err := d.DB().QueryContext(ctx, getSubscribersOfDispatchQ, dispatchId)
	if d.IsError(err, InvalidTextRepresentation) {
		return result, fmt.Errorf("%w: incorrect format for uuid", services.InvalidArgumentErr)
	}
	for rows.Next() {
		var email sql.NullString
		if err := rows.Scan(&email); err != nil {
			return result, fmt.Errorf("failed to scan row: %w", err)
		}
		if email.Valid {
			result = append(result, email.String)
		}
	}

	return result, nil
}

const getAllDispatchesQ = `
	select cd.u_id, cd.label, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	group by cd.id;
`

func (r *DispatchRepo) GetAllDispatches(ctx context.Context, db DB) ([]services.DispatchData, error) {
	result := make([]services.DispatchData, 0, 5) // nolint:mnd
	rows, err := db.DB().QueryContext(ctx, getAllDispatchesQ)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var dispatch services.DispatchData
		if err := rows.Scan(&dispatch.Id, &dispatch.Label, &dispatch.SendAt, &dispatch.CountOfSubscribers); err != nil {
			return result, fmt.Errorf("failed to scan currency dispatch: %w", err)
		}
		result = append(result, dispatch)
	}

	return result, nil
}
