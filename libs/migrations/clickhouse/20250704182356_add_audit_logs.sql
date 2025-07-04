-- +goose Up
CREATE TABLE audit_logs (
	table_name LowCardinality(String),
	operation_type LowCardinality(String),
	old_value Nullable(String),
	new_value Nullable(String),
	object_id Nullable(String),
	user_id LowCardinality(String),
	channel_id LowCardinality(String),
	created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (channel_id, table_name, operation_type, user_id, created_at);
-- +goose Down

