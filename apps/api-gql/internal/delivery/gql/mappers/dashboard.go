package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func DashboardStatsEntityToGql(e entity.DashboardStats) gqlmodel.DashboardStats {
	platforms := make([]gqlmodel.DashboardPlatformStats, 0, len(e.Platforms))
	for _, platformStats := range e.Platforms {
		platforms = append(platforms, gqlmodel.DashboardPlatformStats{
			Platform:     dashboardPlatformToGql(platformStats.Platform),
			IsLive:       platformStats.IsLive,
			Title:        platformStats.Title,
			CategoryID:   platformStats.CategoryID,
			CategoryName: platformStats.CategoryName,
			Viewers:      platformStats.Viewers,
			Followers:    platformStats.Followers,
			StartedAt:    platformStats.StartedAt,
			ChatMessages: platformStats.ChatMessages,
			UsedEmotes:   platformStats.UsedEmotes,
			CanEditInfo:  platformStats.CanEditInfo,
		})
	}

	return gqlmodel.DashboardStats{
		CategoryID:     e.StreamCategoryID,
		CategoryName:   e.StreamCategoryName,
		Viewers:        e.StreamViewers,
		StartedAt:      e.StreamStartedAt,
		Title:          e.StreamTitle,
		ChatMessages:   e.StreamChatMessages,
		Followers:      e.Followers,
		UsedEmotes:     e.UsedEmotes,
		RequestedSongs: e.RequestedSongs,
		Subs:           e.Subs,
		Platforms:      platforms,
	}
}

func dashboardPlatformToGql(platform platformentity.Platform) gqlmodel.Platform {
	switch platform {
	case platformentity.PlatformTwitch:
		return gqlmodel.PlatformTwitch
	case platformentity.PlatformKick:
		return gqlmodel.PlatformKick
	case platformentity.PlatformVKVideoLive:
		return gqlmodel.PlatformVkVideoLive
	default:
		return gqlmodel.Platform(platform)
	}
}

func DashboardBotInfoEntityToGql(e entity.BotStatus) gqlmodel.BotStatus {
	return gqlmodel.BotStatus{
		DashboardID: e.DashboardID,
		Platform:    e.Platform,
		ChannelName: e.ChannelName,
		IsMod:       e.IsMod,
		BotID:       e.BotID,
		BotName:     e.BotName,
		Enabled:     e.Enabled,
	}
}
