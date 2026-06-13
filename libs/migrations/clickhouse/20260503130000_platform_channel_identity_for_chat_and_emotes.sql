-- +goose Up
RENAME TABLE channels_emotes_usages TO channels_emotes_usages_old_platform_identity;

CREATE TABLE channels_emotes_usages
(
	platform LowCardinality(String) DEFAULT 'twitch',
	platform_channel_id LowCardinality(String),
	user_id String,
	created_at DateTime DEFAULT now(),
	emote LowCardinality(String)
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (platform, platform_channel_id, emote, user_id, created_at);

INSERT INTO channels_emotes_usages (platform, platform_channel_id, user_id, created_at, emote)
SELECT 'twitch', channel_id, user_id, created_at, emote
FROM channels_emotes_usages_old_platform_identity;


RENAME TABLE chat_messages TO chat_messages_old_platform_identity;

CREATE TABLE chat_messages
(
	id UUID DEFAULT generateUUIDv4(),
	platform LowCardinality(String) DEFAULT 'twitch',
	platform_channel_id LowCardinality(String),
	user_id LowCardinality(String),
	user_name LowCardinality(String),
	user_display_name LowCardinality(String),
	user_color LowCardinality(String),
	text String,
	created_at DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (platform, platform_channel_id, user_id, created_at, id);

INSERT INTO chat_messages (id, platform, platform_channel_id, user_id, user_name, user_display_name, user_color, text, created_at)
SELECT id, 'twitch', channel_id, user_id, user_name, user_display_name, user_color, text, created_at
FROM chat_messages_old_platform_identity;

-- +goose Down
