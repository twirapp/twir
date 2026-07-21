-- +goose Up
-- +goose StatementBegin
CREATE TABLE channel_platforms (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id          UUID NOT NULL REFERENCES channels(id),
    platform            TEXT NOT NULL,
    user_id             UUID NOT NULL REFERENCES users(id),
    platform_channel_id TEXT NOT NULL,
    enabled             BOOLEAN NOT NULL DEFAULT TRUE,
    bot_user_id         UUID REFERENCES users(id),
    bot_config          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT channel_platforms_channel_id_platform_key UNIQUE (channel_id, platform)
);

CREATE UNIQUE INDEX channel_platforms_platform_channel_id_idx
    ON channel_platforms USING btree (platform, platform_channel_id);

CREATE INDEX channel_platforms_channel_id_idx ON channel_platforms USING btree (channel_id);
CREATE INDEX channel_platforms_user_id_idx ON channel_platforms USING btree (user_id);
CREATE INDEX channel_platforms_bot_user_id_idx ON channel_platforms USING btree (bot_user_id);

INSERT INTO channel_platforms (
    channel_id,
    platform,
    user_id,
    platform_channel_id,
    enabled,
    bot_user_id,
    bot_config
)
SELECT
    c.id,
    'twitch',
    c.twitch_user_id,
    u.platform_id,
    c.twitch_bot_enabled,
    NULL,
    jsonb_build_object(
        'bot_id', c."botId",
        'is_bot_mod', c."isBotMod",
        'is_twitch_banned', c."isTwitchBanned"
    )
FROM channels c
JOIN users u ON u.id = c.twitch_user_id AND u.platform = 'twitch'
WHERE c.twitch_user_id IS NOT NULL
ON CONFLICT DO NOTHING;

INSERT INTO channel_platforms (
    channel_id,
    platform,
    user_id,
    platform_channel_id,
    enabled,
    bot_user_id,
    bot_config
)
SELECT
    c.id,
    'kick',
    c.kick_user_id,
    u.platform_id,
    c.kick_bot_enabled,
    kb.kick_user_id,
    jsonb_strip_nulls(jsonb_build_object('kick_bot_id', c.kick_bot_id))
FROM channels c
JOIN users u ON u.id = c.kick_user_id AND u.platform = 'kick'
LEFT JOIN kick_bots kb ON kb.id = c.kick_bot_id
WHERE c.kick_user_id IS NOT NULL
ON CONFLICT DO NOTHING;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM channels c
        LEFT JOIN users u ON u.id = c.twitch_user_id AND u.platform = 'twitch'
        WHERE c.twitch_user_id IS NOT NULL
          AND (
              u.id IS NULL
              OR NOT EXISTS (
                  SELECT 1
                  FROM channel_platforms cp
                  WHERE cp.channel_id = c.id
                    AND cp.platform = 'twitch'
                    AND cp.user_id = c.twitch_user_id
                    AND cp.platform_channel_id = u.platform_id
                    AND cp.enabled = c.twitch_bot_enabled
                    AND cp.bot_user_id IS NULL
                    AND cp.bot_config = jsonb_build_object(
                        'bot_id', c."botId",
                        'is_bot_mod', c."isBotMod",
                        'is_twitch_banned', c."isTwitchBanned"
                    )
              )
          )
    ) THEN
        RAISE EXCEPTION 'channel_platforms Twitch backfill did not preserve every legacy channel binding';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM channels c
        LEFT JOIN users u ON u.id = c.kick_user_id AND u.platform = 'kick'
        LEFT JOIN kick_bots kb ON kb.id = c.kick_bot_id
        WHERE c.kick_user_id IS NOT NULL
          AND (
              u.id IS NULL
              OR (c.kick_bot_id IS NOT NULL AND kb.id IS NULL)
              OR NOT EXISTS (
                  SELECT 1
                  FROM channel_platforms cp
                  WHERE cp.channel_id = c.id
                    AND cp.platform = 'kick'
                    AND cp.user_id = c.kick_user_id
                    AND cp.platform_channel_id = u.platform_id
                    AND cp.enabled = c.kick_bot_enabled
                    AND cp.bot_user_id IS NOT DISTINCT FROM kb.kick_user_id
                    AND cp.bot_config = jsonb_strip_nulls(jsonb_build_object('kick_bot_id', c.kick_bot_id))
              )
          )
    ) THEN
        RAISE EXCEPTION 'channel_platforms Kick backfill did not preserve every legacy channel binding';
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channel_platforms;
-- +goose StatementEnd
