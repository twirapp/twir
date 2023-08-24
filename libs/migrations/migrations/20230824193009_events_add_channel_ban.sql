-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_events_type_enum ADD VALUE 'CHANNEL_BAN';
ALTER TYPE channel_events_list_type_enum ADD VALUE 'CHANNEL_BAN';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
