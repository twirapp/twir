-- +goose Up
-- +goose StatementBegin

-- channels_greetings."userId" was never remapped during the users multi-platform
-- migration (20260415001328_users_multi_platform.sql), because the column has no
-- FK to users. Legacy rows still contain raw Twitch platform IDs instead of
-- internal users UUIDs, which breaks pgx uuid scanning in the greetings repository.

-- 1. Drop legacy rows that would collide with already-remapped/new rows
--    (unique index channels_greetings_userid_channelid_unique_idx on ("userId", "channelId"))
DELETE
FROM channels_greetings g
WHERE g."userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
  AND EXISTS (
      SELECT 1
      FROM users u
               JOIN channels_greetings g2
                    ON g2."channelId" = g."channelId" AND g2."userId" = u.id::text
      WHERE u.platform = 'twitch'
        AND u.platform_id = g."userId"
  );

-- 2. Remap legacy Twitch platform IDs to internal user UUIDs
UPDATE channels_greetings g
SET "userId" = u.id::text
FROM users u
WHERE u.platform = 'twitch'
  AND u.platform_id = g."userId"
  AND g."userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

-- 3. Drop orphaned rows whose user no longer exists
DELETE
FROM channels_greetings
WHERE "userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

-- 4. Convert the column to a real uuid and rename it to user_id
ALTER TABLE channels_greetings
    ALTER COLUMN "userId" TYPE uuid USING "userId"::uuid;

ALTER TABLE channels_greetings
    RENAME COLUMN "userId" TO user_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE channels_greetings
    RENAME COLUMN user_id TO "userId";

ALTER TABLE channels_greetings
    ALTER COLUMN "userId" TYPE text USING "userId"::text;

-- +goose StatementEnd
