-- +goose Up
CREATE TABLE IF NOT EXISTS file_processing_completed_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payload JSONB NOT NULL,
    status SMALLINT NOT NULL,
    attempts SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_file_processing_completed_outbox_status_created
ON file_processing_completed_outbox (status, created_at);

CREATE INDEX IF NOT EXISTS idx_file_processing_completed_outbox_status_updated_attempts_created
ON file_processing_completed_outbox (status, updated_at, attempts, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_file_processing_completed_outbox_status_created;
DROP INDEX IF EXISTS idx_file_processing_completed_outbox_status_updated_attempts_created;
DROP TABLE IF EXISTS file_processing_completed_outbox;
