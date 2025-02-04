alter table question
    drop constraint if exists question_author_id_fkey,
    drop column if exists author_id;