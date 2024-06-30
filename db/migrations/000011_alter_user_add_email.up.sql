alter table "user" add column  if not exists email text;
alter table "user" add constraint unique_user_email unique (email);