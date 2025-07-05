-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS channels_streams_channel_id_index ON channels_streams("userId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
