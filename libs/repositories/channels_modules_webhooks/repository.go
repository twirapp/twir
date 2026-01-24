package channels_modules_webhooks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/webhook_notifications"
)

var ErrSettingsNotFound = errors.New("webhook module settings not found")

type Repository interface {
	GetByID(ctx context.Context, id string) (webhook_notifications.Settings, error)
	GetByChannelID(ctx context.Context, channelID string) (webhook_notifications.Settings, error)
	GetEnabledChannels(ctx context.Context, input GetEnabledChannelsInput) ([]string, error)
	Create(ctx context.Context, input CreateInput) (webhook_notifications.Settings, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (webhook_notifications.Settings, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID string

	Enabled                   bool
	GithubIssuesEnabled       bool
	GithubPullRequestsEnabled bool
	GithubCommitsEnabled      bool
}

type UpdateInput struct {
	Enabled                   *bool
	GithubIssuesEnabled       *bool
	GithubPullRequestsEnabled *bool
	GithubCommitsEnabled      *bool
}

type GetEnabledChannelsInput struct {
	GithubIssuesEnabled       *bool
	GithubPullRequestsEnabled *bool
	GithubCommitsEnabled      *bool
	Page                      int
	PerPage                   int
}
