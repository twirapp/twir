-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels
    DROP CONSTRAINT IF EXISTS "FK_4f890144c0cb55fe7867b8f61e6",
    DROP CONSTRAINT IF EXISTS channels_kick_bot_id_fkey,
    DROP CONSTRAINT IF EXISTS channels_kick_user_id_fkey,
    DROP CONSTRAINT IF EXISTS channels_twitch_user_id_fkey;

DROP INDEX IF EXISTS channels_kick_bot_enabled_idx;
DROP INDEX IF EXISTS channels_kick_bot_id_idx;
DROP INDEX IF EXISTS channels_kick_user_id_unique_idx;
DROP INDEX IF EXISTS channels_twitch_bot_enabled_idx;
DROP INDEX IF EXISTS channels_twitch_user_id_unique_idx;

ALTER TABLE channels
    DROP COLUMN IF EXISTS twitch_user_id,
    DROP COLUMN IF EXISTS kick_user_id,
    DROP COLUMN IF EXISTS kick_bot_id,
    DROP COLUMN IF EXISTS "botId",
    DROP COLUMN IF EXISTS "isBotMod",
    DROP COLUMN IF EXISTS "isTwitchBanned",
    DROP COLUMN IF EXISTS twitch_bot_enabled,
    DROP COLUMN IF EXISTS kick_bot_enabled;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels
    ADD COLUMN IF NOT EXISTS twitch_user_id uuid,
    ADD COLUMN IF NOT EXISTS kick_user_id uuid,
    ADD COLUMN IF NOT EXISTS kick_bot_id uuid,
    ADD COLUMN IF NOT EXISTS "botId" text,
    ADD COLUMN IF NOT EXISTS "isBotMod" boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS "isTwitchBanned" boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS twitch_bot_enabled boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS kick_bot_enabled boolean NOT NULL DEFAULT false;

UPDATE channels c
SET
    twitch_user_id = cp.user_id,
    twitch_bot_enabled = cp.enabled,
    "isBotMod" = COALESCE((cp.bot_config ->> 'is_bot_mod')::boolean, false),
    "isTwitchBanned" = COALESCE((cp.bot_config ->> 'is_twitch_banned')::boolean, false),
    "botId" = cp.bot_config ->> 'bot_id'
FROM channel_platforms cp
WHERE cp.channel_id = c.id
    AND cp.platform = 'twitch';

UPDATE channels c
SET
    kick_user_id = cp.user_id,
    kick_bot_enabled = cp.enabled,
    kick_bot_id = (cp.bot_config ->> 'kick_bot_id')::uuid
FROM channel_platforms cp
WHERE cp.channel_id = c.id
    AND cp.platform = 'kick';

ALTER TABLE channels
    ADD CONSTRAINT channels_twitch_user_id_fkey FOREIGN KEY (twitch_user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    ADD CONSTRAINT channels_kick_user_id_fkey FOREIGN KEY (kick_user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    ADD CONSTRAINT channels_kick_bot_id_fkey FOREIGN KEY (kick_bot_id) REFERENCES kick_bots (id) ON UPDATE CASCADE ON DELETE SET NULL,
    ADD CONSTRAINT "FK_4f890144c0cb55fe7867b8f61e6" FOREIGN KEY ("botId") REFERENCES bots (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE UNIQUE INDEX channels_twitch_user_id_unique_idx ON channels (twitch_user_id);
CREATE UNIQUE INDEX channels_kick_user_id_unique_idx ON channels (kick_user_id);
CREATE INDEX channels_twitch_bot_enabled_idx ON channels (twitch_bot_enabled);
CREATE INDEX channels_kick_bot_enabled_idx ON channels (kick_bot_enabled);
CREATE INDEX channels_kick_bot_id_idx ON channels (kick_bot_id);
-- +goose StatementEnd
