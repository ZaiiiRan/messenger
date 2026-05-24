-- +goose Up
CREATE TABLE IF NOT EXISTS file_variants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    variant SMALLINT NOT NULL,
    bucket SMALLINT NOT NULL,
    s3_key TEXT NOT NULL,
    mime_type TEXT,
    size BIGINT NOT NULL DEFAULT 0,
    meta JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_file_variants_file_id ON file_variants (file_id);

CREATE TYPE v1_file_variant AS (
    id UUID,
    file_id UUID,
    variant SMALLINT,
    bucket SMALLINT,
    s3_key TEXT,
    mime_type TEXT,
    size BIGINT,
    meta JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_file_variant;
DROP INDEX IF EXISTS idx_file_variants_file_id;
DROP TABLE IF EXISTS file_variants;
