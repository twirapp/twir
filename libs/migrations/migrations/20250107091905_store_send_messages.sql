-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE sent_messages (
	id SERIAL PRIMARY KEY,
	twitch_id TEXT NOT NULL,
	content TEXT,
	channel_id TEXT,
	sender_id TEXT,
	created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
