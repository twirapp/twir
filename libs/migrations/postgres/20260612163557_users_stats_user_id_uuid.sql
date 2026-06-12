-- +goose Up
-- +goose StatementBegin
-- Drop old indexes that reference the old column names
DROP INDEX IF EXISTS users_stats_userid_idx;
DROP INDEX IF EXISTS "users_stats_userId_channelId_key";

-- Delete orphaned rows where userId doesn't match any users.id
DELETE FROM users_stats
WHERE "userId" NOT IN (SELECT id::text FROM users);

-- Rename columns: camelCase -> snake_case
ALTER TABLE users_stats RENAME COLUMN "userId" TO user_id;
ALTER TABLE users_stats RENAME COLUMN "channelId" TO channel_id;

-- Convert TEXT columns to UUID
ALTER TABLE users_stats ALTER COLUMN user_id TYPE UUID USING user_id::uuid;
ALTER TABLE users_stats ALTER COLUMN channel_id TYPE UUID USING channel_id::uuid;

-- Add foreign key constraints
ALTER TABLE users_stats
    ADD CONSTRAINT users_stats_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE users_stats
    ADD CONSTRAINT users_stats_channel_id_fkey
    FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;

-- Recreate unique index with new column names
CREATE UNIQUE INDEX users_stats_user_id_channel_id_key ON users_stats(user_id, channel_id);
CREATE INDEX users_stats_user_id_idx ON users_stats(user_id);

-- Fix channels_streams.userId: the users_multi_platform migration stored
-- internal UUIDs instead of Twitch platform IDs. Restore correct values.
UPDATE channels_streams cs
SET "userId" = u.platform_id
FROM users u
WHERE u.id::text = cs."userId"
  AND u.platform = 'twitch'
  AND u.platform_id != cs."userId";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'users_stats_user_id_uuid is not reversible';
-- +goose StatementEnd
