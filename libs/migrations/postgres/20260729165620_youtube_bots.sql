-- +goose Up
-- +goose StatementBegin
CREATE TABLE youtube_bots (
    id                      UUID        PRIMARY KEY DEFAULT uuidv7(),
    singleton               BOOLEAN     NOT NULL DEFAULT TRUE,
    encrypted_access_token  TEXT        NOT NULL,
    encrypted_refresh_token TEXT        NOT NULL,
    scopes                  TEXT[]      NOT NULL DEFAULT '{}',
    expires_in              INTEGER     NOT NULL DEFAULT 0,
    obtainment_timestamp    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    youtube_user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT youtube_bots_singleton_true CHECK (singleton),
    CONSTRAINT youtube_bots_singleton_key UNIQUE (singleton)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS youtube_bots;
-- +goose StatementEnd
