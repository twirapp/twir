-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_events_type_enum ADD VALUE 'CHANNEL_MESSAGE_DELETE';
ALTER TYPE channels_events_operations_type_enum ADD VALUE 'MESSAGE_DELETE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
