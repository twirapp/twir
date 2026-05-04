-- +goose Up
-- +goose StatementBegin
CREATE TABLE kick_bots (
    id                   UUID        PRIMARY KEY DEFAULT uuidv7(),
    type                 TEXT        NOT NULL DEFAULT 'bot',
    access_token         TEXT        NOT NULL,
    refresh_token        TEXT        NOT NULL,
    scopes               TEXT[]      NOT NULL DEFAULT '{}',
    expires_in           INTEGER     NOT NULL DEFAULT 0,
    obtainment_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    kick_user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    kick_user_login      TEXT        NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX kick_bots_kick_user_id_key ON kick_bots(kick_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kick_bots;
-- +goose StatementEnd
