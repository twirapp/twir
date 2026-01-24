package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/entities/webhook_notifications"
)

func WebhookNotificationsEntityToGql(
	e webhook_notifications.Settings,
) gqlmodel.WebhookNotificationsSettings {
	return gqlmodel.WebhookNotificationsSettings{
		ID:                 e.ID,
		ChannelID:          e.ChannelID,
		Enabled:            e.Enabled,
		GithubIssues:       e.GithubIssuesEnabled,
		GithubPullRequests: e.GithubPullRequestsEnabled,
		GithubCommits:      e.GithubCommitsEnabled,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}
