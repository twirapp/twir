-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE "channels_messages";
DROP INDEX IF EXISTS "PK_685840c6efdbd345cf265976753";
DROP INDEX IF EXISTS "IDX_05a589d7e48f170714dc73243b";
DROP INDEX IF EXISTS "IDX_dd5560724c1166b9b70954f44b";

CREATE UNLOGGED TABLE "channels_messages" (
	"id" SERIAL PRIMARY KEY,
	"message_id" TEXT,
	"channel_id" TEXT references channels(id) ON DELETE CASCADE,
	"user_id" TEXT,
	"user_name" TEXT,
	"text" TEXT,
	"can_be_deleted" BOOLEAN,
	"created_at" TIMESTAMP
);

CREATE INDEX "channels_messages_channel_id_index" ON "channels_messages" ("channel_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
