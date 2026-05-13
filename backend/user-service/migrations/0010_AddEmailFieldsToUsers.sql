-- +goose Up
ALTER TABLE status
    ADD COLUMN old_email TEXT,
    ADD COLUMN email_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

DROP TYPE IF EXISTS v1_status;
CREATE TYPE v1_status AS (
    id BIGINT,
    user_id UUID,
    is_confirmed BOOLEAN,
    is_permanently_banned BOOLEAN,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN,
    deleted_at TIMESTAMPTZ,
    is_permanently_deleted BOOLEAN,
    old_email TEXT,
    email_updated_at TIMESTAMPTZ
);

-- +goose Down
ALTER TABLE status
    DROP COLUMN old_email,
    DROP COLUMN email_updated_at;

DROP TYPE IF EXISTS v1_status;
CREATE TYPE v1_status AS (
    id BIGINT,
    user_id UUID,
    is_confirmed BOOLEAN,
    is_permanently_banned BOOLEAN,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN,
    deleted_at TIMESTAMPTZ,
    is_permanently_deleted BOOLEAN
);
