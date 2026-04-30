alter table if exists users add column if not exists uid UUID not null default gen_random_uuid() UNIQUE;

create index idx_users_uid  on users(uid);