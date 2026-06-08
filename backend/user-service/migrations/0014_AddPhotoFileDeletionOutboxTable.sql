-- +goose Up
CREATE TABLE photo_file_deletion_outbox (
    id          BIGSERIAL    PRIMARY KEY,
    payload     JSONB        NOT NULL,
    status      SMALLINT     NOT NULL DEFAULT 0,
    attempts    SMALLINT     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_photo_file_deletion_outbox_status_created_at
    ON photo_file_deletion_outbox(status, created_at);
CREATE INDEX idx_photo_file_deletion_outbox_status_updated_at_attempts
    ON photo_file_deletion_outbox(status, updated_at, attempts, created_at);

-- +goose Down
DROP TABLE IF EXISTS photo_file_deletion_outbox;
