-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE channels_giveaways DROP COLUMN "is_running";
ALTER TABLE channels_giveaways DROP COLUMN "is_finished";
ALTER TABLE channels_giveaways DROP COLUMN "is_stopped";
ALTER TABLE channels_giveaways DROP COLUMN "is_archived";

ALTER TABLE channels_giveaways ADD COLUMN "stopped_at" TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
