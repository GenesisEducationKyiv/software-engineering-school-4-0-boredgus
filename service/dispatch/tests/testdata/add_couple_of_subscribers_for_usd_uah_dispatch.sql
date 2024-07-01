insert into subs."users" (email)
values ('test_email_1@gmail.com'), ('test_email_2@gmail.com');

insert into subs."currency_subscriptions" (user_id, dispatch_id)
select *
from 
  (select u.id user_id
  from subs."users" as u
  where u.email = 'test_email_1@gmail.com' or u.email = 'test_email_2@gmail.com') u,
  (select cd.id dispatch_id
  from subs.currency_dispatches cd
  where cd.u_id = 'f669a90d-d4aa-4285-bbce-6b14c6ff9065') d;
