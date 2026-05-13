-- +goose Up
CREATE TABLE IF NOT EXISTS change_email_tokens (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            UUID NOT NULL,
    email              TEXT NOT NULL,
    code               TEXT NOT NULL,
    link_token         TEXT NOT NULL,
    generations_left   INTEGER NOT NULL,
    verifications_left INTEGER NOT NULL,
    expires_at         TIMESTAMPTZ NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_change_email_tokens_user_id ON change_email_tokens(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_change_email_tokens_link_token ON change_email_tokens(link_token);

-- +goose Down
DROP INDEX IF EXISTS idx_change_email_tokens_link_token;
DROP INDEX IF EXISTS idx_change_email_tokens_user_id;
DROP TABLE IF EXISTS change_email_tokens;

