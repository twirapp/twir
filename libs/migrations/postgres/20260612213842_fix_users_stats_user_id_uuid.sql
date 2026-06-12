-- +goose NO TRANSACTION
-- +goose Up

-- STEP 1: Convert old TEXT user_id values to UUID strings via platform_id mapping
-- +goose StatementBegin
UPDATE users_stats us
SET user_id = u.id::text
FROM users u
WHERE u.platform = 'twitch'
	AND u.platform_id = us.user_id
	AND us.user_id !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';
-- +goose StatementEnd

-- STEP 2: Drop indexes BEFORE type conversion
-- +goose StatementBegin
DROP INDEX IF EXISTS users_stats_user_id_channel_id_key;
DROP INDEX IF EXISTS users_stats_user_id_idx;
-- +goose StatementEnd

-- STEP 3: Convert user_id column from TEXT to UUID
-- +goose StatementBegin
ALTER TABLE users_stats
	ALTER COLUMN user_id TYPE UUID USING user_id::uuid;
-- +goose StatementEnd

-- STEP 4: Deduplicate — delete non-canonical rows first, then update survivors
-- +goose StatementBegin
DELETE
FROM users_stats
WHERE id IN (SELECT id
             FROM (SELECT id,
                          ROW_NUMBER() OVER (
														PARTITION BY user_id, channel_id
														ORDER BY created_at ASC
														) AS rn
                   FROM users_stats) sub
             WHERE rn > 1);
-- +goose StatementEnd

-- +goose StatementBegin
WITH ranked AS (SELECT id,
                       SUM(messages) OVER (PARTITION BY user_id, channel_id)            AS total_messages,
                       SUM(watched) OVER (PARTITION BY user_id, channel_id)             AS total_watched,
                       SUM(emotes) OVER (PARTITION BY user_id, channel_id)              AS total_emotes,
                       SUM("usedChannelPoints")
                       OVER (PARTITION BY user_id, channel_id)                          AS total_pts,
                       SUM(reputation) OVER (PARTITION BY user_id, channel_id)          AS total_rep,
                       BOOL_OR(is_mod) OVER (PARTITION BY user_id, channel_id)          AS any_mod,
                       BOOL_OR(is_vip) OVER (PARTITION BY user_id, channel_id)          AS any_vip,
                       BOOL_OR(is_subscriber) OVER (PARTITION BY user_id, channel_id)   AS any_sub,
                       MIN(created_at) OVER (PARTITION BY user_id, channel_id)          AS earliest_created,
                       MAX(updated_at) OVER (PARTITION BY user_id, channel_id)          AS latest_updated
                FROM users_stats)
UPDATE users_stats us
SET messages            = r.total_messages,
    watched             = r.total_watched,
    emotes              = r.total_emotes,
    "usedChannelPoints" = r.total_pts,
    reputation          = r.total_rep,
    is_mod              = r.any_mod,
    is_vip              = r.any_vip,
    is_subscriber       = r.any_sub,
    created_at          = r.earliest_created,
    updated_at          = r.latest_updated
FROM ranked r
WHERE us.id = r.id;
-- +goose StatementEnd

-- STEP 5: Recreate indexes and add FK
-- +goose StatementBegin
CREATE UNIQUE INDEX users_stats_user_id_channel_id_key ON users_stats (user_id, channel_id);
CREATE INDEX users_stats_user_id_idx ON users_stats (user_id);
ALTER TABLE users_stats
	ADD CONSTRAINT users_stats_user_id_fkey
		FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;
-- +goose StatementEnd

-- STEP 6: Fix users_online.userId
-- +goose StatementBegin
UPDATE users_online uo
SET "userId" = u.id::text
FROM users u
WHERE u.platform = 'twitch'
	AND u.platform_id = uo."userId"
	AND uo."userId" !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users_online
	ALTER COLUMN "userId" TYPE UUID USING "userId"::uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'fix_users_stats_user_id_uuid is not reversible';
-- +goose StatementEnd
