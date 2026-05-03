-- +goose Up
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_email;

CREATE INDEX IF NOT EXISTS idx_user_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_email    ON users(email);

-- +goose Down
DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_email;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email    ON users(email);

ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT users_email_key    UNIQUE (email);
