-- +goose Up
CREATE TABLE IF NOT EXISTS confirmation_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    code TEXT NOT NULL,
    generations_left INTEGER NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_confirmation_codes_user_id ON confirmation_codes(user_id);

CREATE TYPE v1_confirmation_code AS (
    id BIGINT,
    user_id UUID,
    code TEXT,
    generations_left INTEGER,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_confirmation_codes_user_id;
DROP TYPE IF EXISTS v1_confirmation_code;
DROP TABLE IF EXISTS confirmation_codes;
