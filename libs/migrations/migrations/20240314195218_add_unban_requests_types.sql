-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channel_events_list_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_CREATE';
ALTER TYPE channel_events_list_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_RESOLVE';

ALTER TYPE channels_events_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_CREATE';
ALTER TYPE channels_events_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_RESOLVE';

ALTER TYPE channels_modules_settings_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_CREATE';
ALTER TYPE channels_modules_settings_type_enum ADD VALUE 'CHANNEL_UNBAN_REQUEST_RESOLVE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
