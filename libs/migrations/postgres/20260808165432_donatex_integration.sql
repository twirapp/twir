-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels_integrations_donatex (
	id SERIAL PRIMARY KEY,
	public_id UUID NOT NULL DEFAULT uuidv7(),
	channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
	access_token TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	donatex_user_id TEXT NOT NULL DEFAULT '',
	username TEXT NOT NULL,
	avatar TEXT NOT NULL,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE(channel_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_integrations_donatex;
-- +goose StatementEnd
