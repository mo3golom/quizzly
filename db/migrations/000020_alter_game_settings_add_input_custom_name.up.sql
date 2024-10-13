alter table game_settings
    add column if not exists input_custom_name boolean not null default false;