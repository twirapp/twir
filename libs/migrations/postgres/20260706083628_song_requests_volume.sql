-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_song_requests_settings
	ADD COLUMN volume INT NOT NULL DEFAULT 30;

ALTER TABLE channels_song_requests_settings
	ADD CONSTRAINT channels_song_requests_settings_volume_check CHECK (volume >= 0 AND volume <= 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_song_requests_settings
	DROP CONSTRAINT IF EXISTS channels_song_requests_settings_volume_check;

ALTER TABLE channels_song_requests_settings
	DROP COLUMN IF EXISTS volume;
-- +goose StatementEnd
