-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX channels_events_list_channel_id_index ON channels_events_list(channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
