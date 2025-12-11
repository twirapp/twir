-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS channels_integrations_faceit
(
	id           SERIAL PRIMARY KEY,
	channel_id   TEXT        NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
	access_token TEXT        NOT NULL,
	username     TEXT        NOT NULL,
	avatar       TEXT        NOT NULL,
	game         TEXT        NOT NULL DEFAULT '',
	enabled      BOOLEAN     NOT NULL DEFAULT TRUE,
	created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (channel_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS channels_integrations_faceit;
-- +goose StatementEnd
