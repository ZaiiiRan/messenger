CREATE DATABASE messenger;
USE messenger;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    phone TEXT UNIQUE,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    birthdate DATE,

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    is_activated BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    refresh_token TEXT NOT NULL
);

CREATE TABLE friend_statuses (
    id SMALLSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO friend_statuses (name) VALUES ('request'), ('accepted'), ('blocked');

CREATE TABLE friends (
    friend_1_id BIGINT NOT NULL REFERENCES users(id),
    friend_2_id BIGINT NOT NULL REFERENCES users(id),
    status_id SMALLINT NOT NULL REFERENCES friend_statuses(id)
);

CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    name TEXT

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE chat_roles (
    id SMALLSERIAL PRIMARY KEY,
    role TEXT NOT NULL UNIQUE
);

INSERT INTO chat_roles (role) VALUES ('member'), ('admin'), ('owner');

CREATE TABLE chat_members (
    chat_id BIGINT NOT NULL REFERENCES chats(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    role_id SMALLINT NOT NULL,

    removed_by BIGINT REFERENCES users(id) -- NULL if user left the chat
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_edited TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE messages_read_status (
    message_id BIGINT NOT NULL REFERENCES messages(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP
);

CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(id),
    file_url TEXT,
    file_size INTEGER
);

CREATE TABLE activation_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);


-- Trigger for checking members count in chats
-- (maybe it will be useful in future)

CREATE OR REPLACE FUNCTION check_chat_members_limit()
    RETURNS TRIGGER AS $$
    DECLARE 
        chat_type BOOLEAN;
        member_count INTEGER;
    BEGIN
        SELECT is_group_chat INTO chat_type FROM chats WHERE id = NEW.chat_id;

        IF NOT chat_type THEN
            SELECT COUNT(*) INTO member_count FROM chat_members WHERE chat_id = NEW.chat_id AND is_deleted = FALSE;
    
            IF member_count >= 2 THEN
                RAISE EXCEPTION 'This chat is not a group chat and cannot have more than 2 members';
            END IF;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER chat_members_limit_trigger
BEFORE INSERT ON chat_members
FOR EACH ROW
EXECUTE FUNCTION check_chat_members_limit();

