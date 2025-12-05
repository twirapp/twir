-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_modules_obs_websocket
	ADD COLUMN IF NOT EXISTS scenes        TEXT[] NOT NULL DEFAULT '{}',
	ADD COLUMN IF NOT EXISTS sources       TEXT[] NOT NULL DEFAULT '{}',
	ADD COLUMN IF NOT EXISTS audio_sources TEXT[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_modules_obs_websocket
	DROP COLUMN IF EXISTS scenes,
	DROP COLUMN IF EXISTS sources,
	DROP COLUMN IF EXISTS audio_sources;
-- +goose StatementEnd

