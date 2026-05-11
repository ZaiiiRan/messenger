-- +goose Up
CREATE TYPE v1_outbox_event AS (
    id UUID,
    payload JSONB,
    status SMALLINT,
    attempts SMALLINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TYPE IF EXISTS v1_outbox_event;
