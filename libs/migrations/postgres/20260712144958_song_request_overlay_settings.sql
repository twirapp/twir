-- +goose Up
-- +goose StatementBegin
CREATE TYPE song_request_overlay_style AS ENUM (
	'CINEMA',
	'COMPACT',
	'TICKER',
	'STUDIO',
	'PORTRAIT',
	'PILL'
);

CREATE TABLE channels_song_requests_overlay_settings (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id UUID NOT NULL UNIQUE REFERENCES channels(id) ON DELETE CASCADE,
	style song_request_overlay_style NOT NULL DEFAULT 'CINEMA',
	accent_color VARCHAR(9) NOT NULL DEFAULT '#8B5CF6',
	ticker_background_color VARCHAR(9) NOT NULL DEFAULT '#111827E6',
	ticker_text_color VARCHAR(9) NOT NULL DEFAULT '#FFFFFF',
	ticker_speed INT NOT NULL DEFAULT 35,
	hide_on_pause BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT song_request_overlay_accent_color_check
		CHECK (accent_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$'),
	CONSTRAINT song_request_overlay_ticker_background_color_check
		CHECK (ticker_background_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$'),
	CONSTRAINT song_request_overlay_ticker_text_color_check
		CHECK (ticker_text_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$'),
	CONSTRAINT song_request_overlay_ticker_speed_check
		CHECK (ticker_speed BETWEEN 10 AND 100)
);

INSERT INTO channels_song_requests_overlay_settings (channel_id, hide_on_pause)
SELECT channel_id, COALESCE(hide_on_pause, true)
FROM channels_song_requests_settings
ON CONFLICT (channel_id) DO NOTHING;

ALTER TABLE channels_song_requests_settings
	DROP COLUMN hide_on_pause;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_song_requests_settings
	ADD COLUMN hide_on_pause BOOLEAN DEFAULT true;

UPDATE channels_song_requests_settings AS settings
SET hide_on_pause = overlay.hide_on_pause
FROM channels_song_requests_overlay_settings AS overlay
WHERE overlay.channel_id = settings.channel_id;

DROP TABLE channels_song_requests_overlay_settings;
DROP TYPE song_request_overlay_style;
-- +goose StatementEnd
