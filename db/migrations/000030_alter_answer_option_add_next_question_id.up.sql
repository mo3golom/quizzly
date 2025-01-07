alter table question_answer_option
    add column if not exists next_question_id uuid references question(id) default null;