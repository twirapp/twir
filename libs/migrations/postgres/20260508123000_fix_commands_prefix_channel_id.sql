-- +goose Up
-- +goose StatementBegin
UPDATE channels_commands_prefix ccp
SET channel_id = c.id::text
FROM users u
JOIN channels c ON c.twitch_user_id = u.id
WHERE u.platform = 'twitch'
  AND u.platform_id = ccp.channel_id;

WITH ranked AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY channel_id
            ORDER BY updated_at DESC NULLS LAST, created_at DESC NULLS LAST, id DESC
        ) AS rn
    FROM channels_commands_prefix
)
DELETE FROM channels_commands_prefix ccp
USING ranked r
WHERE ccp.id = r.id
  AND r.rn > 1;

CREATE UNIQUE INDEX IF NOT EXISTS channels_commands_prefix_channel_id_unique_idx
    ON channels_commands_prefix (channel_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_commands_prefix_channel_id_unique_idx;
-- +goose StatementEnd
