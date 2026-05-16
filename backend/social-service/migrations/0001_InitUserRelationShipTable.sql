-- +goose Up
CREATE TABLE IF NOT EXISTS user_relationships (
    user_id_1 UUID NOT NULL,
    user_id_2 UUID NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id_1, user_id_2)
);

CREATE INDEX IF NOT EXISTS idx_user_relationships_user2
ON user_relationships (user_id_2);

CREATE TYPE v1_user_relationship AS (
    user_id_1 UUID,
    user_id_2 UUID,
    status SMALLINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP INDEX IF EXISTS idx_user_relationships_user2;
DROP TYPE IF EXISTS v1_user_relationship;
DROP TABLE IF EXISTS user_relationships;
