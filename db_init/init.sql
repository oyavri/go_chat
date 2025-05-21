DROP TABLE IF EXISTS chat_user;
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS chat_message;
DROP TABLE IF EXISTS chat_member;

CREATE TABLE chat_user (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN NOT NULL
);

CREATE TABLE chat (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE chat_message (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    chat_id uuid NOT NULL,
    content VARCHAR,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chat (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES chat_user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE chat_member (
    chat_id uuid NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chat (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES chat_user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
