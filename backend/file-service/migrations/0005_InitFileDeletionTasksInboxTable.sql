-- +goose Up
CREATE TABLE IF NOT EXISTS file_deletion_tasks_inbox (
    id UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    status SMALLINT NOT NULL,
    attempts SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_file_deletion_tasks_inbox_status_created
ON file_deletion_tasks_inbox (status, created_at);

CREATE INDEX IF NOT EXISTS idx_file_deletion_tasks_inbox_status_updated_attempts_created
ON file_deletion_tasks_inbox (status, updated_at, attempts, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_file_deletion_tasks_inbox_status_created;
DROP INDEX IF EXISTS idx_file_deletion_tasks_inbox_status_updated_attempts_created;
DROP TABLE IF EXISTS file_deletion_tasks_inbox;
