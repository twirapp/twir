-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS channels_categories_aliases (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		channel_id TEXT NOT NULL,
		alias TEXT NOT NULL,
		category_id TEXT NOT NULL,

		FOREIGN KEY (channel_id) REFERENCES channels (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE channels_categories_aliases;
-- +goose StatementEnd
