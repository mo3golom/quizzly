create table if not exists question_answer_option (
    id bigint generated by default as identity primary key not null,
    question_id UUID not null,
    answer text not null,
    is_correct boolean not null,
    created_at TIMESTAMPTZ not null default NOW(),
    updated_at TIMESTAMPTZ not null default NOW(),

    foreign key (question_id) references question (id)
)