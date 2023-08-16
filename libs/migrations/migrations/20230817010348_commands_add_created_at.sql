-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_commands_usages"
	ADD COLUMN "createdAt" timestamp default now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "channels_commands_usages"
	DROP COLUMN "createdAt"
-- +goose StatementEnd
