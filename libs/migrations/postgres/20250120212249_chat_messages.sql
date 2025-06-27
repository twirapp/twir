-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- ID        uuid.UUID
-- 	ChannelID string
-- 	UserID    string
-- 	Text      string
-- 	CreatedAt time.Time
-- 	UpdateAt  time.Time
CREATE UNLOGGED TABLE chat_messages (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	channel_id TEXT REFERENCES channels(id) ON DELETE CASCADE NOT NULL,
	user_id TEXT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
	user_name TEXT NOT NULL,
	user_display_name TEXT NOT NULL,
	user_color TEXT NOT NULL,
	text TEXT NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
