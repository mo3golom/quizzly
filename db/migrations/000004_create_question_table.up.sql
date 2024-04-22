create table if not exists question (
    id UUID primary key not null,
    "text" text not null,
    "type" text not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW()
)