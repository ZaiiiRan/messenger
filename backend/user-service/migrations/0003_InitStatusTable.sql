-- +goose Up
CREATE TABLE IF NOT EXISTS status (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    is_permanently_banned BOOLEAN NOT NULL DEFAULT FALSE,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_status_user_id ON status(user_id);

CREATE TYPE v1_status AS (
    id BIGINT,
    user_id UUID,
    is_confirmed BOOLEAN,
    is_permanently_banned BOOLEAN,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN
);

-- +goose Down
DROP INDEX IF EXISTS idx_status_user_id;
DROP TYPE IF EXISTS v1_status;
DROP TABLE IF EXISTS status;
