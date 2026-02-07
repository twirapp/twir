-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_giveaways_settings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    winner_message TEXT NOT NULL DEFAULT 'Congratulations {winner}! You won the giveaway!',
    FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);

CREATE INDEX idx_channels_giveaways_settings_channel_id ON channels_giveaways_settings(channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_giveaways_settings;
-- +goose StatementEnd
