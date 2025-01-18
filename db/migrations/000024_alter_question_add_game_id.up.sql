alter table question drop constraint if exists game_game_id;
alter table question
    add column if not exists game_id UUID,
    add constraint game_game_id foreign key (game_id) references game(id);