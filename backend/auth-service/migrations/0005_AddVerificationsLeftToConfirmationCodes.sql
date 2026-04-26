-- +goose Up
ALTER TABLE confirmation_codes
    ADD COLUMN verifications_left INTEGER NOT NULL DEFAULT 10;

DROP TYPE IF EXISTS v1_confirmation_code;
CREATE TYPE v1_confirmation_code AS (
    id               BIGINT,
    user_id          UUID,
    code             TEXT,
    generations_left INTEGER,
    verifications_left INTEGER,
    expires_at       TIMESTAMPTZ,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ
);

-- +goose Down
ALTER TABLE confirmation_codes DROP COLUMN IF EXISTS verifications_left;

DROP TYPE IF EXISTS v1_confirmation_code;
CREATE TYPE v1_confirmation_code AS (
    id               BIGINT,
    user_id          UUID,
    code             TEXT,
    generations_left INTEGER,
    expires_at       TIMESTAMPTZ,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ
);
