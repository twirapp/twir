package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapDiscordIntegrationDataToGql(data *entity.DiscordIntegrationData) *gqlmodel.DiscordIntegrationData {
	if data == nil {
		return &gqlmodel.DiscordIntegrationData{
			Guilds: []gqlmodel.DiscordGuild{},
		}
	}

	guilds := make([]gqlmodel.DiscordGuild, 0, len(data.Guilds))
	for _, guild := range data.Guilds {
		guilds = append(guilds, MapDiscordGuildToGql(guild))
	}

	return &gqlmodel.DiscordIntegrationData{
		Guilds: guilds,
	}
}

func MapDiscordGuildToGql(guild entity.DiscordGuild) gqlmodel.DiscordGuild {
	return gqlmodel.DiscordGuild{
		ID:                               guild.ID,
		Name:                             guild.Name,
		Icon:                             guild.Icon,
		LiveNotificationEnabled:          guild.LiveNotificationEnabled,
		LiveNotificationChannelsIds:      guild.LiveNotificationChannelsIds,
		LiveNotificationShowTitle:        guild.LiveNotificationShowTitle,
		LiveNotificationShowCategory:     guild.LiveNotificationShowCategory,
		LiveNotificationShowViewers:      guild.LiveNotificationShowViewers,
		LiveNotificationMessage:          guild.LiveNotificationMessage,
		LiveNotificationShowPreview:      guild.LiveNotificationShowPreview,
		LiveNotificationShowProfileImage: guild.LiveNotificationShowProfileImage,
		OfflineNotificationMessage:       guild.OfflineNotificationMessage,
		ShouldDeleteMessageOnOffline:     guild.ShouldDeleteMessageOnOffline,
		AdditionalUsersIdsForLiveCheck:   guild.AdditionalUsersIdsForLiveCheck,
	}
}

func MapDiscordChannelTypeToGql(t entity.DiscordChannelType) gqlmodel.DiscordChannelType {
	switch t {
	case entity.DiscordChannelTypeText:
		return gqlmodel.DiscordChannelTypeChannelTypeText
	case entity.DiscordChannelTypeVoice:
		return gqlmodel.DiscordChannelTypeChannelTypeVoice
	case entity.DiscordChannelTypeDM:
		return gqlmodel.DiscordChannelTypeChannelTypeDm
	case entity.DiscordChannelTypeGroupDM:
		return gqlmodel.DiscordChannelTypeChannelTypeGroupDm
	case entity.DiscordChannelTypeCategory:
		return gqlmodel.DiscordChannelTypeChannelTypeCategory
	case entity.DiscordChannelTypeAnnouncement:
		return gqlmodel.DiscordChannelTypeChannelTypeAnnouncement
	case entity.DiscordChannelTypeAnnouncementThread:
		return gqlmodel.DiscordChannelTypeChannelTypeAnnouncementThread
	case entity.DiscordChannelTypePublicThread:
		return gqlmodel.DiscordChannelTypeChannelTypePublicThread
	case entity.DiscordChannelTypePrivateThread:
		return gqlmodel.DiscordChannelTypeChannelTypePrivateThread
	case entity.DiscordChannelTypeStageVoice:
		return gqlmodel.DiscordChannelTypeChannelTypeStageVoice
	case entity.DiscordChannelTypeDirectory:
		return gqlmodel.DiscordChannelTypeChannelTypeDirectory
	case entity.DiscordChannelTypeForum:
		return gqlmodel.DiscordChannelTypeChannelTypeForum
	case entity.DiscordChannelTypeMedia:
		return gqlmodel.DiscordChannelTypeChannelTypeMedia
	default:
		return gqlmodel.DiscordChannelTypeChannelTypeText
	}
}

func MapDiscordGuildChannelsToGql(channels []entity.DiscordGuildChannel) []gqlmodel.DiscordGuildChannel {
	result := make([]gqlmodel.DiscordGuildChannel, 0, len(channels))
	for _, channel := range channels {
		result = append(result, gqlmodel.DiscordGuildChannel{
			ID:              channel.ID,
			Name:            channel.Name,
			Type:            MapDiscordChannelTypeToGql(channel.Type),
			CanSendMessages: channel.CanSendMessages,
		})
	}
	return result
}

func MapDiscordGuildInfoToGql(info *entity.DiscordGuildInfo) *gqlmodel.DiscordGuildInfo {
	if info == nil {
		return nil
	}

	channels := make([]gqlmodel.DiscordGuildChannel, 0, len(info.Channels))
	for _, channel := range info.Channels {
		channels = append(channels, gqlmodel.DiscordGuildChannel{
			ID:              channel.ID,
			Name:            channel.Name,
			Type:            MapDiscordChannelTypeToGql(channel.Type),
			CanSendMessages: channel.CanSendMessages,
		})
	}

	roles := make([]gqlmodel.DiscordGuildRole, 0, len(info.Roles))
	for _, role := range info.Roles {
		roles = append(roles, gqlmodel.DiscordGuildRole{
			ID:    role.ID,
			Name:  role.Name,
			Color: role.Color,
		})
	}

	return &gqlmodel.DiscordGuildInfo{
		ID:       info.ID,
		Name:     info.Name,
		Icon:     info.Icon,
		Channels: channels,
		Roles:    roles,
	}
}
