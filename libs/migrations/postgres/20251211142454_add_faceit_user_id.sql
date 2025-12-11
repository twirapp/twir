-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_integrations_faceit
	ADD COLUMN IF NOT EXISTS faceit_user_id TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE channels_integrations_faceit
	DROP COLUMN IF EXISTS faceit_user_id;
-- +goose StatementEnd

