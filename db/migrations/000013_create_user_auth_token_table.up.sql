create table if not exists user_auth_token (
    user_id    UUID primary key not null,
    token      text not null,
    expires_at TIMESTAMPTZ not null default NOW(),

    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (user_id) references user (id),
)