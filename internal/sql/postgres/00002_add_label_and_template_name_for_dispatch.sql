-- +goose Up
-- +goose StatementBegin
alter table subs."currency_dispatches"
add column label varchar(60),
add column template_name varchar(60);
-- +goose StatementEnd

-- +goose StatementBegin
update subs."currency_dispatches" cd
set label = 'Exchange rate from USD to UAH',
  template_name = 'exchange_rate'
where cd.u_id = 'f669a90d-d4aa-4285-bbce-6b14c6ff9065';
-- +goose StatementEnd

-- +goose StatementBegin
alter table subs."currency_dispatches"
alter column label set not null,
alter column template_name set not null;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
alter table subs."currency_dispatches"
drop column label,
drop column template_name;
-- +goose StatementEnd
