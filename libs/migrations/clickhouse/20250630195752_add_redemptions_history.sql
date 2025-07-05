-- +goose Up
CREATE TABLE channels_redemptions_history (
	channel_id LowCardinality(String),
	user_id LowCardinality(String),
	reward_id LowCardinality(UUID),
	reward_title LowCardinality(String),
	reward_prompt Nullable(String),
	reward_cost Int32,
	created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (channel_id, reward_id, user_id, created_at);

-- +goose Down

