create table if not exists player (
    id UUID primary key not null,
    user_id UUID default null,
    name text not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (user_id) references "user" (id)
)