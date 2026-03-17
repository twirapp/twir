-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_modules_webhooks
	ADD COLUMN github_issues_online_enabled BOOL NOT NULL DEFAULT true,
	ADD COLUMN github_issues_offline_enabled BOOL NOT NULL DEFAULT true,
	ADD COLUMN github_pull_requests_online_enabled BOOL NOT NULL DEFAULT true,
	ADD COLUMN github_pull_requests_offline_enabled BOOL NOT NULL DEFAULT true,
	ADD COLUMN github_commits_online_enabled BOOL NOT NULL DEFAULT true,
	ADD COLUMN github_commits_offline_enabled BOOL NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_modules_webhooks
	DROP COLUMN github_issues_online_enabled,
	DROP COLUMN github_issues_offline_enabled,
	DROP COLUMN github_pull_requests_online_enabled,
	DROP COLUMN github_pull_requests_offline_enabled,
	DROP COLUMN github_commits_online_enabled,
	DROP COLUMN github_commits_offline_enabled;
-- +goose StatementEnd
