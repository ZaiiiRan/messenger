-- +goose Up
CREATE TABLE IF NOT EXISTS file_processing_tasks (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    payload JSONB NOT NULL,
    status SMALLINT NOT NULL,
    attempts SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_file_processing_tasks_status_created
ON file_processing_tasks (status, created_at);

CREATE INDEX IF NOT EXISTS idx_file_processing_tasks_status_updated_attempts_created
ON file_processing_tasks (status, updated_at, attempts, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_file_processing_tasks_status_created;
DROP INDEX IF EXISTS idx_file_processing_tasks_status_updated_attempts_created;
DROP TABLE IF EXISTS file_processing_tasks;
