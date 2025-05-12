UPDATE user 
SET %s 
WHERE id = $%d 
RETURNING 
    id,
    username,
    email,
    created_at,
    updated_at,
    deleted_at,
    deleted
