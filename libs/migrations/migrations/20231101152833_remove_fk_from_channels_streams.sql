-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channels_streams" DROP CONSTRAINT "FK_d2b9d6113cdeb816207be291ffa";
ALTER TABLE "discord_sended_notifications" DROP CONSTRAINT "discord_sended_notifications_channel_id_fkey";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
