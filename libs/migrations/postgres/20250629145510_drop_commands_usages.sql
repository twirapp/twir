-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE channels_commands_usages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
