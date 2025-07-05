-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels_commands_usages
(
	channel_id String,
	user_id    String,
	command_id UUID,
	created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(created_at)
	ORDER BY (channel_id, command_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
