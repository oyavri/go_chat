INSERT INTO chat_message (user_id, chat_id, content) 
VALUES ($1, $2, $3)
RETURNING 
    id, 
    user_id, 
    chat_id, 
    content, 
    created_at
