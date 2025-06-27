-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE toxic_messages (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		channel_id TEXT REFERENCES channels(id),
		reply_to_user_id TEXT REFERENCES users(id) NULL,
		text TEXT,
		created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
