-- +goose Up
CREATE TABLE IF NOT EXISTS update_friends_in_privacy_settings_inbox (
    id UUID PRIMARY KEY,
    payload JSONB NOT NULL,
    status SMALLINT NOT NULL,
    attempts SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_update_friends_in_privacy_settings_inbox_status_created
ON update_friends_in_privacy_settings_inbox (status, created_at);

CREATE INDEX IF NOT EXISTS idx_update_friends_in_privacy_settings_inbox_status_updated_attempts_created
ON update_friends_in_privacy_settings_inbox (status, updated_at, attempts, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_update_friends_in_privacy_settings_inbox_status_created;
DROP INDEX IF EXISTS idx_update_friends_in_privacy_settings_inbox_status_updated_attempts_created;
DROP TABLE IF EXISTS update_friends_in_privacy_settings_inbox;
