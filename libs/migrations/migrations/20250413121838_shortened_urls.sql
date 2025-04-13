-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE shortened_urls (
	short_id text PRIMARY KEY NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	url TEXT NOT NULL,

	created_by_user_id TEXT,
	FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX shortened_urls_url_idx ON shortened_urls(url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
