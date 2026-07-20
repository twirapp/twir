-- +goose Up
-- +goose StatementBegin
CREATE TABLE dota_match_settlements (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    match_id BIGINT NOT NULL,
    won BOOLEAN NOT NULL,
    mmr_delta INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, match_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dota_match_settlements;
-- +goose StatementEnd
