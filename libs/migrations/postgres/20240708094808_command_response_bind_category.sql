-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_commands_responses" ADD COLUMN "twitch_category_id" text[] NOT NULL default '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
