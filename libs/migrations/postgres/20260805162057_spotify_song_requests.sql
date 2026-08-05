-- +goose Up
-- +goose StatementBegin
CREATE TABLE spotify_song_requests (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id TEXT NOT NULL REFERENCES channels(id),
	track_id TEXT NOT NULL,
	track_uri TEXT NOT NULL,
	title TEXT NOT NULL,
	artist TEXT NOT NULL,
	album TEXT NOT NULL,
	duration_ms INT NOT NULL,
	requester_user_id TEXT,
	requester_name TEXT NOT NULL,
	requester_display_name TEXT,
	source TEXT NOT NULL,
	queue_position INT NOT NULL DEFAULT 0,
	status TEXT NOT NULL,
	queued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	playing_observed_at TIMESTAMPTZ,
	played_observed_at TIMESTAMPTZ,
	cancelled_pending_skip_at TIMESTAMPTZ,
	skipped_by_twir_at TIMESTAMPTZ,
	removed_or_reconciled_at TIMESTAMPTZ,
	unknown_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX spotify_song_requests_channel_id_status_idx
	ON spotify_song_requests (channel_id, status);

CREATE INDEX spotify_song_requests_channel_id_requester_name_status_idx
	ON spotify_song_requests (channel_id, requester_name, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE spotify_song_requests;
-- +goose StatementEnd
