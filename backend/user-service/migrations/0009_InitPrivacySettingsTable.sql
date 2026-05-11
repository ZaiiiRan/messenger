-- +goose Up
CREATE TABLE IF NOT EXISTS privacy_settings (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    settings JSONB NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_privacy_settings_user_id ON privacy_settings(user_id);

CREATE TYPE v1_privacy_settings AS (
    id BIGINT,
    user_id UUID,
    settings JSONB
);

-- +goose Down
DROP INDEX IF EXISTS idx_privacy_settings_user_id;
DROP TYPE IF EXISTS v1_privacy_settings;
DROP TABLE IF EXISTS privacy_settings;
