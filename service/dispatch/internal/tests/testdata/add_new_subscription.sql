insert into subs."users" (email)
values ('created@gmail.com');

insert into subs."currency_subscriptions" (user_id, dispatch_id, status)
select u.user_id, d.dispatch_id, 1
from 
  (select u.id user_id
  from subs."users" as u
  where u.email = 'created@gmail.com') u,
  (select cd.id dispatch_id
  from subs.currency_dispatches cd
  where cd.u_id = 'f669a90d-d4aa-4285-bbce-6b14c6ff9065') d;
