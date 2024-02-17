-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE channel_overlay_now_playing_preset AS ENUM ('TRANSPARENT', 'AIDEN_REDESIGN');

CREATE TABLE channels_overlays_now_playing (
	id UUID PRIMARY KEY default gen_random_uuid() NOT NULL,
	preset channel_overlay_now_playing_preset,
	channel_id TEXT,
	FOREIGN KEY (channel_id) REFERENCES channels(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
