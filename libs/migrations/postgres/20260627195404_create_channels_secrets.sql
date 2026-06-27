-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_secrets (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    value       TEXT NOT NULL,
    "channel_id" UUID NOT NULL REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE("channel_id", name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels_secrets;
-- +goose StatementEnd
