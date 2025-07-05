-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_emotes_usages
(
	id UUID DEFAULT generateUUIDv4(),
	channel_id String,
	user_id String,
	created_at DateTime DEFAULT now(),
	emote String
) ENGINE = MergeTree() PRIMARY KEY (id) ORDER BY (id, channel_id, user_id, emote, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
