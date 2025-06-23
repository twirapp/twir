-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_commands_responses ADD COLUMN online_only boolean NOT NULL DEFAULT false;
ALTER TABLE channels_commands_responses ADD COLUMN offline_only boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
