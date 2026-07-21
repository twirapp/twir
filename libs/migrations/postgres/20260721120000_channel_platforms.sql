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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channel_platforms;
-- +goose StatementEnd
