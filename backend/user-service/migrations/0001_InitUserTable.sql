-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email ON users(email);

CREATE TYPE v1_user AS (
    id UUID,
    username TEXT,
    email TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_email;
DROP TYPE IF EXISTS v1_user;
DROP TABLE IF EXISTS users;
