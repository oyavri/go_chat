INSERT INTO chat_user (username, email) 
VALUES ($1, $2) 
RETURNING  
    id,
    username,
    email,
    created_at,
    updated_at,
    deleted_at,
    deleted
