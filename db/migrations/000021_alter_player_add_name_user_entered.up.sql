alter table player
    add column if not exists name_user_entered boolean not null default false;