-- +goose Up
ALTER TABLE refresh_tokens
    ADD COLUMN ip      TEXT,
    ADD COLUMN country TEXT,
    ADD COLUMN city    TEXT,
    ADD COLUMN os      TEXT,
    ADD COLUMN browser TEXT;

DROP TYPE IF EXISTS v1_refresh_token;
CREATE TYPE v1_refresh_token AS (
    id         BIGINT,
    user_id    UUID,
    token      TEXT,
    version    INTEGER,
    ip         TEXT,
    country    TEXT,
    city       TEXT,
    os         TEXT,
    browser    TEXT,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_refresh_token;
CREATE TYPE v1_refresh_token AS (
    id         BIGINT,
    user_id    UUID,
    token      TEXT,
    version    INTEGER,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

ALTER TABLE refresh_tokens
    DROP COLUMN IF EXISTS ip,
    DROP COLUMN IF EXISTS country,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS os,
    DROP COLUMN IF EXISTS browser;
