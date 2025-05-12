UPDATE user 
SET deleted = true, 
    deleted_at = NOW() 
WHERE id = $1 
RETURNING 
    id,
    username,
    email,
    created_at,
    updated_at,
    deleted_at,
    deleted
