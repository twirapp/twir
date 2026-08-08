-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels_dota_settings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL UNIQUE REFERENCES channels(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT false,
    steam_account_id TEXT,
    gsi_token TEXT NOT NULL DEFAULT replace(uuidv7()::text, '-', ''),
    mmr INT NOT NULL DEFAULT 0,
    mmr_delta INT NOT NULL DEFAULT 25,
    session_wins INT NOT NULL DEFAULT 0,
    session_losses INT NOT NULL DEFAULT 0,
    prediction_settings JSONB NOT NULL DEFAULT '{"enabled": false, "titleTemplate": "Win this game?", "windowSeconds": 300}'::jsonb,
    chat_events JSONB NOT NULL DEFAULT '{}'::jsonb,
    commands_settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS channels_dota_settings_gsi_token_idx ON channels_dota_settings(gsi_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_dota_settings;
-- +goose StatementEnd
