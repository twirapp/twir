-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_overlays_stream_stats (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id          UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    design              channels_overlays_stream_stats_design NOT NULL DEFAULT 'GLASS',
    variant             channels_overlays_stream_stats_variant NOT NULL DEFAULT 'HORIZONTAL',
    viewers_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    viewers_mode        channels_overlays_stream_stats_viewers_mode NOT NULL DEFAULT 'CUMULATIVE',
    platform_icons_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    messages_enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    uptime_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    subscribers_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    followers_enabled   BOOLEAN NOT NULL DEFAULT TRUE,
    viewers_color       TEXT NOT NULL DEFAULT '',
    messages_color      TEXT NOT NULL DEFAULT '',
    uptime_color        TEXT NOT NULL DEFAULT '',
    subscribers_color   TEXT NOT NULL DEFAULT '',
    followers_color     TEXT NOT NULL DEFAULT '',
    counter_order       TEXT[] NOT NULL DEFAULT '{viewers,messages,uptime,subscribers,followers}',
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
