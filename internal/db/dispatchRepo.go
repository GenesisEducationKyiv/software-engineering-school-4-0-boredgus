package db

import (
	"context"
	"fmt"
	"strings"
	e "subscription-api/internal/entities"
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

func (r *DispatchRepo) GetDispatchByID(ctx context.Context, db DB, dispatchId string) (e.CurrencyDispatch, error) {
	var d e.CurrencyDispatch
	row := db.DB().QueryRowContext(ctx, getDispatchByIdQ, dispatchId)
	err := row.Err()
	if err != nil && db.IsError(err, InvalidTextRepresentation) {
		return d, fmt.Errorf("%w: incorrect format for uuid", services.InvalidArgumentErr)
	}
	if err != nil {
		return d, err
	}
	var targetCurrencies string
	if err := row.Scan(&d.Id, &d.Label, &d.TemplateName, &d.Details.BaseCurrency, &targetCurrencies, &d.SendAt, &d.CountOfSubscribers); err != nil {
		return d, fmt.Errorf("%w: dispatch with such id does not exists", services.NotFoundErr)
	}
	d.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")

	return d, nil
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
		var email string
		if err := rows.Scan(&email); err != nil {
			return result, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, email)
	}

	return result, nil
}

const getAllDispatchesQ = `
	select cd.u_id, cd.label, cd.template_name, cd.base_currency , cd.target_currencies, cd.send_at, count(cs.user_id) subs_count
	from subs."currency_dispatches" cd
	left join subs."currency_subscriptions" cs
	on cs.dispatch_id = cd.id
	group by cd.id;
`

func (r *DispatchRepo) GetAllDispatches(ctx context.Context, db DB) ([]e.CurrencyDispatch, error) {
	// dispatchCount := 5
	result := make([]e.CurrencyDispatch, 0, 5)
	rows, err := db.DB().QueryContext(ctx, getAllDispatchesQ)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d e.CurrencyDispatch
		var targetCurrencies string
		if err := rows.Scan(&d.Id, &d.Label, &d.TemplateName, &d.Details.BaseCurrency, &targetCurrencies, &d.SendAt, &d.CountOfSubscribers); err != nil {
			return result, fmt.Errorf("failed to scan currency dispatch: %w", err)
		}
		d.Details.TargetCurrencies = strings.Split(targetCurrencies, ",")
		result = append(result, d)
	}

	return result, nil
}
