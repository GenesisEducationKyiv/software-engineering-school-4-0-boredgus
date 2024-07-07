-- +goose Up
-- +goose StatementBegin
create table subs."subscription_status" (
  id serial,
  label varchar(60) unique,
  constraint id_pk primary key (id)
);
-- +goose StatementEnd


-- +goose StatementBegin
insert into subs."subscription_status" (label)
values ('created'), ('cancelled'), ('renewed');
-- +goose StatementEnd

-- +goose StatementBegin
alter table subs."currency_subscriptions"
add column status int not null default 1,
add column updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
add constraint status_fk 
  foreign key (status) 
  references subs."subscription_status" (id);
-- +goose StatementEnd

-- +goose StatementBegin
create function subs.on_update_current_timestamp() returns trigger
  language plpgsql
  as $$
begin
   NEW.updated_at = now();
   return NEW;
end;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
create trigger on_update_current_timestamp
before update on subs.currency_subscriptions
for each row execute function subs.on_update_current_timestamp();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop trigger on_update_current_timestamp
on table subs.currency_subscriptions;
-- +goose StatementEnd

-- +goose StatementBegin
drop function subs.on_update_current_timestamp;
-- +goose StatementEnd

-- +goose StatementBegin
alter table subs."currency_subscriptions"
drop column status,
drop column updated_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop table subs."subscription_status";
-- +goose StatementEnd
