-- +goose StatementBegin

-- The users_multi_platform migration missed users_stats.user_id because
-- there was no FK constraint from users_stats.user_id to users.id.
-- As a result, users_stats.user_id still contains old Twitch TEXT platform IDs
-- (e.g. "684505240") instead of UUIDs, making 99% of stats invisible to code.
-- Same issue affects users_online.userId.
--
-- This migration:
-- 1. Converts old TEXT user_id to UUID via users.platform_id mapping
-- 2. Merges duplicate rows (old TEXT + new UUID for same user+channel)
-- 3. Converts column types and adds FK constraints

-- ============================================================
-- STEP 1: Convert old TEXT user_id values to UUID
-- Map via users.platform_id (Twitch ID) -> users.id (UUID)
-- ============================================================
UPDATE users_stats us
SET user_id = u.id::text
FROM users u
WHERE u.platform = 'twitch'
  AND u.platform_id = us.user_id
  AND us.user_id !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

-- ============================================================
-- STEP 2: Merge duplicates — sum old values into the earliest row per (user_id, channel_id)
-- After step 1, some user+channel pairs have 2 rows: old (with data) + new (with 0s).
-- ============================================================
WITH ranked AS (
    SELECT
        id,
        user_id,
        channel_id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, channel_id
            ORDER BY created_at ASC
        ) AS rn,
        SUM(messages)            OVER (PARTITION BY user_id, channel_id) AS total_messages,
        SUM(watched)             OVER (PARTITION BY user_id, channel_id) AS total_watched,
        SUM(emotes)              OVER (PARTITION BY user_id, channel_id) AS total_emotes,
        SUM("usedChannelPoints") OVER (PARTITION BY user_id, channel_id) AS total_pts,
        SUM(reputation)          OVER (PARTITION BY user_id, channel_id) AS total_rep,
        -- pick is_mod/vip/sub from the row that has TRUE (any), prefer old (earliest created_at)
        BOOL_OR(is_mod)          OVER (PARTITION BY user_id, channel_id) AS any_mod,
        BOOL_OR(is_vip)          OVER (PARTITION BY user_id, channel_id) AS any_vip,
        BOOL_OR(is_subscriber)   OVER (PARTITION BY user_id, channel_id) AS any_sub,
        MIN(created_at)          OVER (PARTITION BY user_id, channel_id) AS earliest_created,
        MAX(updated_at)          OVER (PARTITION BY user_id, channel_id) AS latest_updated
    FROM users_stats
),
to_update AS (
    SELECT * FROM ranked WHERE rn = 1
),
to_delete AS (
    SELECT id FROM ranked WHERE rn > 1
)
UPDATE users_stats us
SET
    messages            = tu.total_messages,
    watched             = tu.total_watched,
    emotes              = tu.total_emotes,
    "usedChannelPoints" = tu.total_pts,
    reputation          = tu.total_rep,
    is_mod              = tu.any_mod,
    is_vip              = tu.any_vip,
    is_subscriber       = tu.any_sub,
    created_at          = tu.earliest_created,
    updated_at          = tu.latest_updated
FROM to_update tu
WHERE us.id = tu.id
  AND tu.total_messages IS NOT NULL;

-- Delete the duplicate rows (keep only rn=1 per group)
DELETE FROM users_stats
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY user_id, channel_id ORDER BY created_at ASC
        ) AS rn
        FROM users_stats
    ) sub
    WHERE rn > 1
);

-- ============================================================
-- STEP 3: Drop indexes before type change
-- ============================================================
DROP INDEX IF EXISTS users_stats_user_id_channel_id_key;
DROP INDEX IF EXISTS users_stats_user_id_idx;

-- ============================================================
-- STEP 4: Convert user_id column from TEXT to UUID
-- ============================================================
ALTER TABLE users_stats ALTER COLUMN user_id TYPE UUID USING user_id::uuid;

-- ============================================================
-- STEP 5: Recreate indexes and add FK
-- ============================================================
CREATE UNIQUE INDEX users_stats_user_id_channel_id_key ON users_stats(user_id, channel_id);
CREATE INDEX users_stats_user_id_idx ON users_stats(user_id);

ALTER TABLE users_stats
    ADD CONSTRAINT users_stats_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- ============================================================
-- STEP 6: Fix users_online.userId (same issue)
-- No merge needed — users_online is ephemeral.
-- ============================================================
UPDATE users_online uo
SET "userId" = u.id::text
FROM users u
WHERE u.platform = 'twitch'
  AND u.platform_id = uo."userId"
  AND uo."userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';

ALTER TABLE users_online ALTER COLUMN "userId" TYPE UUID USING "userId"::uuid;

-- +goose StatementEnd
