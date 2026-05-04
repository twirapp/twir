package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func DashboardStatsEntityToGql(e entity.DashboardStats) gqlmodel.DashboardStats {
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
