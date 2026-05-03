-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels
    ADD COLUMN twitch_bot_enabled BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN kick_bot_enabled   BOOLEAN NOT NULL DEFAULT false;

UPDATE channels
SET
    twitch_bot_enabled = CASE WHEN twitch_user_id IS NOT NULL THEN "isEnabled" ELSE false END,
    kick_bot_enabled   = CASE WHEN kick_user_id IS NOT NULL THEN "isEnabled" ELSE false END;

CREATE INDEX channels_twitch_bot_enabled_idx ON channels(twitch_bot_enabled);
CREATE INDEX channels_kick_bot_enabled_idx ON channels(kick_bot_enabled);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_twitch_bot_enabled_idx;
DROP INDEX IF EXISTS channels_kick_bot_enabled_idx;

ALTER TABLE channels
    DROP COLUMN IF EXISTS twitch_bot_enabled,
    DROP COLUMN IF EXISTS kick_bot_enabled;
-- +goose StatementEnd
