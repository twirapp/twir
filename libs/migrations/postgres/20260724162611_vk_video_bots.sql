-- +goose Up
-- +goose StatementBegin
CREATE TABLE vk_video_bots (
    id                      UUID        PRIMARY KEY DEFAULT uuidv7(),
    singleton               BOOLEAN     NOT NULL DEFAULT TRUE,
    encrypted_access_token  TEXT        NOT NULL,
    encrypted_refresh_token TEXT        NOT NULL,
    scopes                  TEXT[]      NOT NULL DEFAULT '{}',
    expires_in              INTEGER     NOT NULL DEFAULT 0,
    obtainment_timestamp    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    vk_user_id              UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT vk_video_bots_singleton_true CHECK (singleton),
    CONSTRAINT vk_video_bots_singleton_key UNIQUE (singleton)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vk_video_bots;
-- +goose StatementEnd
