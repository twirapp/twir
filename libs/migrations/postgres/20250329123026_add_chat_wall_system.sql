-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- !chat wall delete BAD PHRASE
-- !chat wall ban BAD PHRASE
-- !chat wall timeout 2d BAD PHRASE

-- !chat wall stop BAD PHRASE

CREATE TYPE channels_chat_wall_action AS ENUM (
	'DELETE',
	'BAN',
	'TIMEOUT'
);

CREATE TABLE IF NOT EXISTS channels_chat_wall_settings (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	mute_subscribers BOOL NOT NULL DEFAULT true,
	mute_vips BOOL NOT NULL DEFAULT false,

	FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX channels_chat_wall_settings_channel_id_unique_idx ON channels_chat_wall_settings (channel_id);

CREATE TABLE IF NOT EXISTS channels_chat_wall (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	phrase TEXT NOT NULL check (length(phrase) > 0 AND length(phrase) <= 1000),
	enabled BOOL NOT NULL,

	action channels_chat_wall_action NOT NULL,
	duration_seconds INT NOT NULL DEFAULT 600,
	timeout_duration_seconds INT,

	FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS channels_chat_wall_log (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	wall_id UUID NOT NULL,
	user_id TEXT NOT NULL,
	text TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	FOREIGN KEY (wall_id) REFERENCES channels_chat_wall (id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX channels_chat_wall_log_wall_id_idx ON channels_chat_wall_log (wall_id);
CREATE INDEX channels_chat_wall_log_user_id_idx ON channels_chat_wall_log (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
