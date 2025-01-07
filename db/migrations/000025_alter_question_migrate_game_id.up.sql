update question
    set game_id = game_question.game_id
    from game_question
    where question.id = game_question.question_id;