-- +goose Up
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            UUID NOT NULL,
    code               TEXT NOT NULL,
    link_token         TEXT NOT NULL,
    generations_left   INTEGER NOT NULL,
    verifications_left INTEGER NOT NULL,
    expires_at         TIMESTAMPTZ NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id    ON password_reset_tokens(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_password_reset_tokens_link_token  ON password_reset_tokens(link_token);

CREATE TYPE v1_password_reset_token AS (
    id                 BIGINT,
    user_id            UUID,
    code               TEXT,
    link_token         TEXT,
    generations_left   INTEGER,
    verifications_left INTEGER,
    expires_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_password_reset_tokens_link_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;
DROP TYPE IF EXISTS v1_password_reset_token;
DROP TABLE IF EXISTS password_reset_tokens;
