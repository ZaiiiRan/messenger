-- +goose Up
CREATE TABLE IF NOT EXISTS profile (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT,
    birthdate DATE,
    bio TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_user_id ON profile(user_id);

CREATE TYPE v1_profile AS (
    id BIGINT,
    user_id UUID,
    first_name TEXT,
    last_name TEXT,
    phone TEXT,
    birthdate DATE,
    bio TEXT
);

-- +goose Down
DROP INDEX IF EXISTS idx_profile_user_id;
DROP TYPE IF EXISTS v1_profile;
DROP TABLE IF EXISTS profile;
