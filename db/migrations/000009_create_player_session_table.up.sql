create table if not exists player_session (
    id bigint generated by default as identity primary key not null,
    game_id UUID not null,
    player_id UUID not null,
    status text not null,

    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (game_id) references game (id),
    foreign key (player_id) references player (id)
)