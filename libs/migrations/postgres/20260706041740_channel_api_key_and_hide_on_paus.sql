-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE channels ADD COLUMN IF NOT EXISTS  api_key TEXT DEFAULT uuidv7();
CREATE UNIQUE INDEX IF NOT EXISTS channels_api_key_idx ON channels(api_key) WHERE api_key IS NOT NULL;

-- Hide on pause setting for song requests
ALTER TABLE channels_song_requests_settings ADD COLUMN hide_on_pause BOOL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
