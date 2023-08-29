-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_alerts"
	ADD COLUMN "greetings_ids" text[] default '{}'::text[];
ALTER TABLE "channels_alerts"
	ADD COLUMN "keywords_ids" text[] default '{}'::text[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
