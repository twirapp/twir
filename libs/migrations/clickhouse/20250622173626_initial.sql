-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_messages
(
	id                UUID     DEFAULT generateUUIDv4(),
	channel_id        String,
	user_id           String,
	user_name         String,
	user_display_name String,
	user_color        String,
	text              String,
	created_at        DateTime DEFAULT now(),
	updated_at        DateTime DEFAULT now()
) ENGINE = MergeTree() ORDER BY (channel_id, created_at, id)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
