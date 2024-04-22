create table if not exists game (
    id UUID primary key not null,
    status text not null,
    "type" text not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW()
)