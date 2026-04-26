-- +goose Up
CREATE TABLE IF NOT EXISTS user_versions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_version_user_id ON user_versions(user_id);

CREATE TYPE v1_user_version AS (
    id BIGINT,
    user_id UUID,
    version INTEGER,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_user_version_user_id;
DROP TYPE IF EXISTS v1_user_version;
DROP TABLE IF EXISTS user_versions;
