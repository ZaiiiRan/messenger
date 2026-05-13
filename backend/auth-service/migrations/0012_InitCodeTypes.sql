-- +goose Up
DROP TYPE IF EXISTS v1_confirmation_code;
DROP TYPE IF EXISTS v1_password_reset_token;

CREATE TYPE v1_code AS (
    id BIGINT,
    user_id UUID,
    code TEXT,
    link_token         TEXT,
    generations_left   INTEGER,
    verifications_left INTEGER,
    expires_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ
);

CREATE TYPE v1_email_code AS (
    id BIGINT,
    user_id UUID,
    email TEXT,
    code TEXT,
    link_token         TEXT,
    generations_left   INTEGER,
    verifications_left INTEGER,
    expires_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_code;
DROP TYPE IF EXISTS v1_email_code;
