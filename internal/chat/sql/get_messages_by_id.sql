SELECT * FROM chat_message 
JOIN chat ON chat_message.chat_id = chat.id 
WHERE chat.id = $1 
ORDER BY chat_message.created_at 
LIMIT $2 
OFFSET $3
