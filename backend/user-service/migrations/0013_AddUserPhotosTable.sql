-- +goose Up
CREATE TABLE user_photos (
    id          UUID         NOT NULL DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id     UUID         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_user_photos_user_id_created_at ON user_photos(user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS user_photos;
