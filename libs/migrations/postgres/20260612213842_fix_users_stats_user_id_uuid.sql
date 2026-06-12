-- +goose Up
-- +goose StatementBegin

-- The users_multi_platform migration missed users_stats.user_id because
-- there was no FK constraint from users_stats.user_id to users.id.
-- As a result, users_stats.user_id still contains old Twitch TEXT platform IDs
-- (e.g. "684505240") instead of UUIDs, making 99% of stats invisible to code.
-- Same issue affects users_online.userId.

-- STEP 1: Convert old TEXT user_id values to UUID strings via platform_id mapping
UPDATE users_stats us
SET user_id = u.id::text
FROM users u
WHERE u.platform = 'twitch'
  AND u.platform_id = us.user_id
  AND us.user_id !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

-- STEP 2: Drop indexes before type change
DROP INDEX IF EXISTS users_stats_user_id_channel_id_key;
DROP INDEX IF EXISTS users_stats_user_id_idx;

-- STEP 3: Convert user_id column from TEXT to UUID
-- This normalizes case (e.g. "019EAE62..." == "019eae62...")
ALTER TABLE users_stats ALTER COLUMN user_id TYPE UUID USING user_id::uuid;

-- STEP 4: Deduplicate by (user_id, channel_id) — keep earliest row, sum values
-- After UUID normalization, different TEXT casings now collide.
DELETE FROM users_stats
WHERE id IN (
    SELECT id FROM (
        SELECT id,
            ROW_NUMBER() OVER (
                PARTITION BY user_id, channel_id
                ORDER BY created_at ASC
            ) AS rn,
            SUM(messages)            OVER (PARTITION BY user_id, channel_id) AS total_messages,
            SUM(watched)             OVER (PARTITION BY user_id, channel_id) AS total_watched,
            SUM(emotes)              OVER (PARTITION BY user_id, channel_id) AS total_emotes,
            SUM("usedChannelPoints") OVER (PARTITION BY user_id, channel_id) AS total_pts,
            SUM(reputation)          OVER (PARTITION BY user_id, channel_id) AS total_rep,
            BOOL_OR(is_mod)          OVER (PARTITION BY user_id, channel_id) AS any_mod,
            BOOL_OR(is_vip)          OVER (PARTITION BY user_id, channel_id) AS any_vip,
            BOOL_OR(is_subscriber)   OVER (PARTITION BY user_id, channel_id) AS any_sub,
            MIN(created_at)          OVER (PARTITION BY user_id, channel_id) AS earliest_created,
            MAX(updated_at)          OVER (PARTITION BY user_id, channel_id) AS latest_updated
        FROM users_stats
    ) sub
    WHERE rn > 1
);

-- Update surviving rows with merged values from deleted duplicates
UPDATE users_stats us
SET
    messages            = merged.total_messages,
    watched             = merged.total_watched,
    emotes              = merged.total_emotes,
    "usedChannelPoints" = merged.total_pts,
    reputation          = merged.total_rep,
    is_mod              = merged.any_mod,
    is_vip              = merged.any_vip,
    is_subscriber       = merged.any_sub,
    created_at          = merged.earliest_created,
    updated_at          = merged.latest_updated
FROM (
    SELECT
        (array_agg(id ORDER BY created_at ASC))[1] AS keep_id,
        SUM(messages)            AS total_messages,
        SUM(watched)             AS total_watched,
        SUM(emotes)              AS total_emotes,
        SUM("usedChannelPoints") AS total_pts,
        SUM(reputation)          AS total_rep,
        BOOL_OR(is_mod)          AS any_mod,
        BOOL_OR(is_vip)          AS any_vip,
        BOOL_OR(is_subscriber)   AS any_sub,
        MIN(created_at)          AS earliest_created,
        MAX(updated_at)          AS latest_updated
    FROM users_stats
    GROUP BY user_id, channel_id
    HAVING COUNT(*) > 1
) merged
WHERE us.id = merged.keep_id;

-- STEP 5: Recreate indexes and add FK
CREATE UNIQUE INDEX users_stats_user_id_channel_id_key ON users_stats(user_id, channel_id);
CREATE INDEX users_stats_user_id_idx ON users_stats(user_id);

ALTER TABLE users_stats
    ADD CONSTRAINT users_stats_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- STEP 6: Fix users_online.userId (same issue, no merge needed)
UPDATE users_online uo
SET "userId" = u.id::text
FROM users u
WHERE u.platform = 'twitch'
  AND u.platform_id = uo."userId"
  AND uo."userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

ALTER TABLE users_online ALTER COLUMN "userId" TYPE UUID USING "userId"::uuid;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'fix_users_stats_user_id_uuid is not reversible';
-- +goose StatementEnd
