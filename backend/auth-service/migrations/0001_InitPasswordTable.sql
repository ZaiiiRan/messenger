-- +goose Up
CREATE TABLE IF NOT EXISTS passwords (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_password_user_id ON passwords(user_id);

CREATE TYPE v1_password AS (
    id BIGINT,
    user_id UUID,
    password_hash TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_password_user_id;
DROP TYPE IF EXISTS v1_password;
DROP TABLE IF EXISTS passwords;
