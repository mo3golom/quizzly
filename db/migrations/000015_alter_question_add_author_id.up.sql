alter table question add column if not exists author_id UUID;
alter table question add foreign key (author_id) references "user" (id);