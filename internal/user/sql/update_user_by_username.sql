UPDATE chat_user 
SET %s 
WHERE username = $%d 
RETURNING 
    id,
    username,
    email,
    created_at,
    updated_at,
    deleted_at,
    deleted
