-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_games_seppuku" (
	"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	"channel_id" TEXT NOT NULL references "channels"("id") ON DELETE CASCADE,
	"enabled" BOOL NOT NULL,
	"timeout_seconds" INT2 NOT NULL,
	"timeout_moderators" BOOL NOT NULL,
	"message" TEXT NOT NULL,
	"message_moderators" TEXT NOT NULL
);

CREATE UNIQUE INDEX "channels_games_seppuku_channel_id_index" ON "channels_games_seppuku" ("channel_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
