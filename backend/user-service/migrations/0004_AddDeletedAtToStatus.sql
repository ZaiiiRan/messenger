-- +goose Up
ALTER TABLE status ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

DROP TYPE IF EXISTS v1_status;
CREATE TYPE v1_status AS (
    id BIGINT,
    user_id UUID,
    is_confirmed BOOLEAN,
    is_permanently_banned BOOLEAN,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN,
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_status;
CREATE TYPE v1_status AS (
    id BIGINT,
    user_id UUID,
    is_confirmed BOOLEAN,
    is_permanently_banned BOOLEAN,
    banned_until TIMESTAMPTZ,
    is_deleted BOOLEAN
);

ALTER TABLE status DROP COLUMN IF EXISTS deleted_at;
