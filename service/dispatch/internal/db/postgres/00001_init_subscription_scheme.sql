-- +goose Up

-- +goose StatementBegin
create schema subs;
-- +goose StatementEnd

-- +goose StatementBegin
alter schema subs owner to postgres;
-- +goose StatementEnd

-- +goose StatementBegin
create table subs."users" (
  id serial,
  email varchar(60) not null,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  constraint user_pk primary key (id),
  constraint unique_email unique (email)
);
-- +goose StatementEnd

-- +goose StatementBegin
create table subs."currency_dispatches" (
  id serial,
  u_id UUID not null default gen_random_uuid(),
  base_currency char(3) not null,
  target_currencies varchar(300) not null,
  send_at time not null,
  constraint dispatch_pk primary key (id),
  constraint unique_id unique (u_id),
  constraint min_one_target_currency check (length(target_currencies) >= 3)
);
-- +goose StatementEnd

-- +goose StatementBegin
create table subs."currency_subscriptions" (
  user_id integer not null,
  dispatch_id integer not null,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  constraint user_fk 
    foreign key (user_id) 
    references subs."users" (id)
    on delete cascade,
  constraint dispatch_fk 
    foreign key (dispatch_id)
    references subs."currency_dispatches" (id)
    on delete cascade,
  constraint unique_subscription unique (user_id, dispatch_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
-- create USD-UAH dispatch
insert into subs."currency_dispatches" (u_id, base_currency, target_currencies, send_at)
values ('f669a90d-d4aa-4285-bbce-6b14c6ff9065','USD','UAH','08:30:00');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP SCHEMA subs CASCADE;
-- +goose StatementEnd
