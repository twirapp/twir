-- +goose Up
RENAME TABLE chat_messages TO chat_messages_old;

CREATE TABLE chat_messages
(
	id                UUID     DEFAULT generateUUIDv4(),
	channel_id        LowCardinality(String),
	user_id           LowCardinality(String),
	user_name         LowCardinality(String),
	user_display_name LowCardinality(String),
	user_color        LowCardinality(String),
	text              String,
	created_at        DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (channel_id, user_id, created_at, id);

INSERT INTO chat_messages (id, channel_id, user_id, user_name, user_display_name, user_color, text, created_at)
SELECT id, channel_id, user_id, user_name, user_display_name, user_color, text, created_at
FROM chat_messages_old;
-- +goose Down
