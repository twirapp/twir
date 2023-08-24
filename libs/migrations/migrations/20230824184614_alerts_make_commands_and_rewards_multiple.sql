-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER table "channels_alerts"
	drop column "command_id";
ALTER table "channels_alerts"
	drop column "reward_id";

ALTER TABLE "channels_alerts"
	ADD COLUMN "command_ids" text[] default '{}'::text[];
ALTER TABLE "channels_alerts"
	ADD COLUMN "reward_ids" text[] default '{}'::text[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
