package channels_modules_webhooks

import "time"

type Settings struct {
	ID                        string    `db:"id"`
	ChannelID                 string    `db:"channel_id"`
	Enabled                   bool      `db:"enabled"`
	GithubIssuesEnabled       bool      `db:"github_issues_enabled"`
	GithubPullRequestsEnabled bool      `db:"github_pull_requests_enabled"`
	GithubCommitsEnabled      bool      `db:"github_commits_enabled"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`

	isNil bool
}

func (s Settings) IsNil() bool {
	return s.isNil
}

var Nil = Settings{isNil: true}
