-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_modules_settings_type_enum ADD VALUE 'chat_alerts';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
