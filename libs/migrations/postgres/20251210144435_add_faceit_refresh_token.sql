-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_integrations_faceit
	ADD COLUMN IF NOT EXISTS refresh_token TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE channels_integrations_faceit
	DROP COLUMN IF EXISTS refresh_token;
-- +goose StatementEnd
