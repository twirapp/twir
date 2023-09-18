-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_commands" ADD COLUMN cooldown_roles_ids text[] default '{}'::text[] NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "channels_commands" DROP COLUMN cooldown_roles_ids;
-- +goose StatementEnd
