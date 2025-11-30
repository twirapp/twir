package channels_integrations_discord

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_integrations_discord/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) ([]model.ChannelIntegrationDiscord, error)
	GetByChannelIDAndGuildID(ctx context.Context, channelID, guildID string) (model.ChannelIntegrationDiscord, error)
	GetByGuildID(ctx context.Context, guildID string) ([]model.ChannelIntegrationDiscord, error)
	GetByAdditionalUserID(ctx context.Context, userID string) ([]model.ChannelIntegrationDiscord, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelIntegrationDiscord, error)
	Update(ctx context.Context, id int, input UpdateInput) error
	Delete(ctx context.Context, id int) error
	DeleteByChannelIDAndGuildID(ctx context.Context, channelID, guildID string) error
}

type CreateInput struct {
	ChannelID                        string
	GuildID                          string
	LiveNotificationEnabled          bool
	LiveNotificationChannelsIds      []string
	LiveNotificationShowTitle        bool
	LiveNotificationShowCategory     bool
	LiveNotificationShowViewers      bool
	LiveNotificationMessage          string
	LiveNotificationShowPreview      bool
	LiveNotificationShowProfileImage bool
	OfflineNotificationMessage       string
	ShouldDeleteMessageOnOffline     bool
	AdditionalUsersIdsForLiveCheck   []string
}

type UpdateInput struct {
	LiveNotificationEnabled          *bool
	LiveNotificationChannelsIds      *[]string
	LiveNotificationShowTitle        *bool
	LiveNotificationShowCategory     *bool
	LiveNotificationShowViewers      *bool
	LiveNotificationMessage          *string
	LiveNotificationShowPreview      *bool
	LiveNotificationShowProfileImage *bool
	OfflineNotificationMessage       *string
	ShouldDeleteMessageOnOffline     *bool
	AdditionalUsersIdsForLiveCheck   *[]string
}
