-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channels_commands_module_enum ADD VALUE 'UTILITY';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
