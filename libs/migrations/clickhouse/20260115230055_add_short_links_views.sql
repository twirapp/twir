-- +goose Up
CREATE TABLE short_links_views (
	short_link_id LowCardinality(String),
	user_id Nullable(String),
	ip Nullable(String),
	user_agent Nullable(String),
	created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (short_link_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS short_links_views;
