alter table game add column if not exists author_id UUID;
alter table game add foreign key (author_id) references "user" (id);