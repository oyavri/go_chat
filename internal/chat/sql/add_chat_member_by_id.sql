INSERT INTO chat_member (chat_id, user_id)
VALUES ($1, $2)
RETURNING
    user_id
