package webhook_notifications

import "time"

type Settings struct {
	ID                        string
	ChannelID                 string
	Enabled                   bool
	GithubIssuesEnabled       bool
	GithubPullRequestsEnabled bool
	GithubCommitsEnabled      bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time

	isNil bool
}

func (s Settings) IsNil() bool {
	return s.isNil
}

var Nil = Settings{
	isNil: true,
}
