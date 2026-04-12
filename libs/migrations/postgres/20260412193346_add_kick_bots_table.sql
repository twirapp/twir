-- +goose Up
-- +goose StatementBegin
CREATE TABLE kick_bots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL DEFAULT 'DEFAULT',
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    scopes TEXT[] NOT NULL DEFAULT '{}',
    expires_in INT NOT NULL DEFAULT 0,
    obtainment_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    kick_user_id TEXT NOT NULL,
    kick_user_login TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE kick_bots;
-- +goose StatementEnd
