-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE channels_modules_webhooks (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	channel_id TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	enabled BOOL NOT NULL DEFAULT false,
	github_issues_enabled BOOL NOT NULL DEFAULT true,
	github_pull_requests_enabled BOOL NOT NULL DEFAULT true,
	github_commits_enabled BOOL NOT NULL DEFAULT true,
	discord_messages_enabled BOOL NOT NULL DEFAULT false,

	FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE channels_modules_webhooks;
-- +goose StatementEnd
