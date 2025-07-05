-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "channels_overlays_dudes" DROP COLUMN "message_box_ignore_commands";

ALTER TABLE "channels_overlays_dudes" ADD COLUMN "ignore_commands" boolean NOT NULL DEFAULT true;
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "ignore_users" boolean NOT NULL DEFAULT true;
ALTER TABLE "channels_overlays_dudes" ADD COLUMN "ignored_users" text[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
