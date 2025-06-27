-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "channel_duels"
	ADD COLUMN "available_until" timestamp null;
ALTER TABLE "channel_duels"
	ADD COLUMN "finished_at" timestamp;
ALTER TABLE "channel_duels"
	ADD COLUMN "sender_login" text default '';
ALTER TABLE "channel_duels"
	ADD COLUMN "sender_moderator" boolean default false;
ALTER TABLE "channel_duels"
	ADD COLUMN "target_login" text default '';
ALTER TABLE "channel_duels"
	ADD COLUMN "target_moderator" boolean default false;

UPDATE "channel_duels"
SET "available_until" = "created_at"
WHERE "available_until" IS NULL;


UPDATE "channel_duels"
SET "finished_at" = "created_at"
WHERE "finished_at" IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
