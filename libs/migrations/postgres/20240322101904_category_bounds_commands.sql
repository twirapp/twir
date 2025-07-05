-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_commands" ADD COLUMN "enabled_categories" text[] default '{}'::text[] NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
