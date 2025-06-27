-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE "channels_modules_settings_type_enum" ADD VALUE 'duel';

CREATE TABLE IF NOT EXISTS "channel_duels" (
	"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
	"channel_id" text NOT NULL references "channels" ("id") ON DELETE CASCADE,
	"sender_id" text references "users" ("id") ON DELETE SET NULL,
	"target_id" text references "users" ("id") ON DELETE SET NULL,
	"loser_id" text references "users" ("id") ON DELETE SET NULL,
	"created_at" timestamp NOT NULL DEFAULT now(),
	PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "channel_duels_channel_id_idx" ON "channel_duels" ("channel_id");
CREATE INDEX IF NOT EXISTS "channel_duels_sender_id_idx" ON "channel_duels" ("sender_id");
CREATE INDEX IF NOT EXISTS "channel_duels_target_id_idx" ON "channel_duels" ("target_id");
CREATE INDEX IF NOT EXISTS "channel_duels_loser_id_idx" ON "channel_duels" ("loser_id");

ALTER TABLE "users_stats" ADD COLUMN "reputation" int8 NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
