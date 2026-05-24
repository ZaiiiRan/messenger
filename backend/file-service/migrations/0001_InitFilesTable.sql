-- +goose Up
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_files_owner_id ON files (owner_id);
CREATE INDEX IF NOT EXISTS idx_files_status ON files (status);

CREATE TYPE v1_file AS (
    id UUID,
    owner_id UUID,
    status SMALLINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_file;
DROP INDEX IF EXISTS idx_files_owner_id;
DROP TABLE IF EXISTS files;
