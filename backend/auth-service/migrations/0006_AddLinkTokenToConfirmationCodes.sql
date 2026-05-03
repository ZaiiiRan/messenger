-- +goose Up
ALTER TABLE confirmation_codes
    ADD COLUMN link_token TEXT;

UPDATE confirmation_codes SET link_token = gen_random_uuid()::text;

ALTER TABLE confirmation_codes
    ALTER COLUMN link_token SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_confirmation_codes_link_token ON confirmation_codes(link_token);

DROP TYPE IF EXISTS v1_confirmation_code;
CREATE TYPE v1_confirmation_code AS (
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
DROP INDEX IF EXISTS idx_confirmation_codes_link_token;
ALTER TABLE confirmation_codes DROP COLUMN IF EXISTS link_token;

DROP TYPE IF EXISTS v1_confirmation_code;
CREATE TYPE v1_confirmation_code AS (
    id                 BIGINT,
    user_id            UUID,
    code               TEXT,
    generations_left   INTEGER,
    verifications_left INTEGER,
    expires_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ
);
