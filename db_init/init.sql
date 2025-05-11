DROP DATABASE IF EXISTS chat;

CREATE DATABASE chat;

USE chat;

CREATE TABLE user (
    id uuid DEFAULT gen_random_uuid(),
    username VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE chat (
    id uuid DEFAULT gen_random_uuid(),
    PRIMARY KEY (id)
);

CREATE TABLE chat_message (
    id uuid DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    chat_id uuid NOT NULL,
    content VARCHAR,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES CHAT (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USER (id)
        ON DELETE CASCADE,
        ON UPDATE CASCADE
);

CREATE TABLE chat_member (
    chat_id uuid NOT NULL,
    user_id uuid NOT NULL
    FOREIGN KEY (chat_id) REFERENCES CHAT (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USER (id)
        ON DELETE CASCADE,
        ON UPDATE CASCADE
);
