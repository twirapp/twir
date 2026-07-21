-- +goose Up
-- +goose StatementBegin
DELETE from channels_streams;

ALTER TABLE channels_streams ADD COLUMN IF NOT EXISTS "platform" platform DEFAULT 'twitch' NOT NULL;
ALTER TABLE channels_streams ADD COLUMN IF NOT EXISTS "channel_id" uuid REFERENCES channels(id) NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS channels_streams_userid_platform ON channels_streams("userId", platform);
CREATE INDEX IF NOT EXISTS channels_streams_channel_id ON channels_streams("channel_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channels_streams_channel_id;
DROP INDEX IF EXISTS channels_streams_userid_platform;

ALTER TABLE channels_streams DROP COLUMN IF EXISTS "channel_id";
ALTER TABLE channels_streams DROP COLUMN IF EXISTS "platform";
-- +goose StatementEnd
