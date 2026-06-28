-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_storage (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX channels_storage_channel_id_key_idx ON channels_storage(channel_id, key);
CREATE INDEX channels_storage_channel_id_idx ON channels_storage(channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_storage;
-- +goose StatementEnd
