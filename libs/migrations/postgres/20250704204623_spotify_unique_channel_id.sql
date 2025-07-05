-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP INDEX IF EXISTS channels_integrations_spotify_channel_id_idx;
CREATE UNIQUE INDEX channels_integrations_spotify_channel_id_idx ON channels_integrations_spotify(channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
