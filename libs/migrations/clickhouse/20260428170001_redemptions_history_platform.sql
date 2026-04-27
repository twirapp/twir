-- +goose Up
ALTER TABLE channels_redemptions_history
	ADD COLUMN IF NOT EXISTS platform LowCardinality(String) DEFAULT 'twitch' AFTER user_id;

-- +goose Down
ALTER TABLE channels_redemptions_history
	DROP COLUMN IF EXISTS platform;
