-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_modules_webhooks
	DROP COLUMN IF EXISTS discord_messages_enabled;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE channels_modules_webhooks
	ADD COLUMN discord_messages_enabled BOOL NOT NULL DEFAULT false;
-- +goose StatementEnd
