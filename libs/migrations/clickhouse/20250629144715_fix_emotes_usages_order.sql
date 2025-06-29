-- +goose Up
RENAME TABLE channels_emotes_usages TO channels_emotes_usages_old;

CREATE TABLE channels_emotes_usages
(
	channel_id String,
	user_id String,
	created_at DateTime DEFAULT now(),
	emote LowCardinality(String)
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (channel_id, emote, user_id, created_at);

INSERT INTO channels_emotes_usages (channel_id, user_id, created_at, emote)
SELECT channel_id, user_id, created_at, emote
FROM channels_emotes_usages_old;

-- DROP TABLE channels_emotes_usages_old;

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
