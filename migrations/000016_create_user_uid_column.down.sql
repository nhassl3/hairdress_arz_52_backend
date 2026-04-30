alter table if exists users drop column if exists uid;

drop index if exists idx_users_uid;