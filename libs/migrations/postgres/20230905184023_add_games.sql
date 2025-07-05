-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE channels_commands_module_enum ADD VALUE 'GAMES';
ALTER TYPE channels_modules_settings_type_enum ADD VALUE '8ball';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
