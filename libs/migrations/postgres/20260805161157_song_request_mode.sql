-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_song_requests_settings
	ADD COLUMN mode VARCHAR(20) NOT NULL DEFAULT 'YOUTUBE';

ALTER TABLE channels_song_requests_settings
	ADD CONSTRAINT channels_song_requests_settings_mode_check CHECK (mode IN ('YOUTUBE', 'SPOTIFY'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_song_requests_settings
	DROP CONSTRAINT IF EXISTS channels_song_requests_settings_mode_check;

ALTER TABLE channels_song_requests_settings
	DROP COLUMN IF EXISTS mode;
-- +goose StatementEnd
