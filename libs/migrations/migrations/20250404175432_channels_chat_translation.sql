-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE channels_chat_translation_settings (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	channel_id TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	enabled BOOL NOT NULL DEFAULT true,
	target_language TEXT NOT NULL,
	excluded_languages TEXT[] NOT NULL DEFAULT '{}',
	use_italic BOOL NOT NULL DEFAULT true,
	excluded_users_ids TEXT[] NOT NULL DEFAULT '{}',

	FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE channels_chat_translation_settings;
-- +goose StatementEnd
