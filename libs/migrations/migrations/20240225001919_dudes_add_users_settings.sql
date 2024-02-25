-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE channels_overlays_dudes_user_settings
(
	id         UUID PRIMARY KEY   default uuid_generate_v4(),
	channel_id TEXT      NOT NULL,
	user_id    TEXT      NOT NULL,
	dude_color TEXT,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	FOREIGN KEY (channel_id) REFERENCES channels (id),
	FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX channels_overlays_dudes_user_settings_channel_id_user_id_idx ON channels_overlays_dudes_user_settings (channel_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
