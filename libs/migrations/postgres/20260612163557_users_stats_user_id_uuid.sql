-- +goose Up
-- +goose StatementBegin
-- After users_multi_platform, users_stats columns are already UUID with FK.
-- Just rename to snake_case and recreate indexes.

-- Drop indexes that reference old camelCase column names
DROP INDEX IF EXISTS users_stats_userid_idx;
DROP INDEX IF EXISTS "users_stats_userId_channelId_key";

-- Rename columns: camelCase -> snake_case
-- FK constraints auto-update their column references on rename
ALTER TABLE users_stats RENAME COLUMN "userId" TO user_id;
ALTER TABLE users_stats RENAME COLUMN "channelId" TO channel_id;

-- Recreate unique index with new column names
CREATE UNIQUE INDEX users_stats_user_id_channel_id_key ON users_stats(user_id, channel_id);
CREATE INDEX users_stats_user_id_idx ON users_stats(user_id);

-- Fix channels_streams.userId: the users_multi_platform migration stored
-- internal UUIDs instead of Twitch platform IDs. Restore correct values.
-- Only touch rows that are valid UUIDs (skip old Twitch platform IDs like "167160215").
UPDATE channels_streams cs
SET "userId" = u.platform_id
FROM users u
WHERE cs."userId" ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
  AND u.id = cs."userId"::uuid
  AND u.platform = 'twitch'
  AND u.platform_id != cs."userId";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'users_stats_user_id_uuid is not reversible';
-- +goose StatementEnd
