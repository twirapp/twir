-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_overlays_stream_stats (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id          UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    design              TEXT NOT NULL DEFAULT 'BAR',
    viewers_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    viewers_mode        TEXT NOT NULL DEFAULT 'CUMULATIVE',
    messages_enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    uptime_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    subscribers_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    followers_enabled   BOOLEAN NOT NULL DEFAULT TRUE,
    custom_html_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    custom_html         TEXT NOT NULL DEFAULT '',
    custom_css          TEXT NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT channels_overlays_stream_stats_channel_id_key UNIQUE (channel_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels_overlays_stream_stats;
-- +goose StatementEnd
