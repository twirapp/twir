-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER table "channels_commands"
	drop column "alert_id";
ALTER TABLE "channels_alerts"
	ADD COLUMN "command_id" uuid null
		constraint "channels_alerts_command_id" references channels_commands ("id")
			on update cascade
			on delete set null;

ALTER TABLE "channels_alerts"
	ADD COLUMN "reward_id" text null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
