create table if not exists user_auth_login_code (
    user_id    UUID primary key not null,
    code       bigint not null,
    expires_at TIMESTAMPTZ not null default NOW(),

    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (user_id) references "user" (id)
)