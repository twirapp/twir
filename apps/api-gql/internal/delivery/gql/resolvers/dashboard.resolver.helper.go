package resolvers

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/redis_keys"
)

func (r *Resolver) getDashboardStats(ctx context.Context) (*gqlmodel.DashboardStats, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, fmt.Errorf("get selected dashboard: %w", err)
	}

	var stream model.ChannelsStreams

	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(
			`"userId" = ?`,
			dashboardId,
		).
		Find(&stream).Error; err != nil {
		return nil, fmt.Errorf("get stream: %w", err)
	}

	result := &gqlmodel.DashboardStats{
		Viewers:      &stream.ViewerCount,
		CategoryID:   stream.GameId,
		CategoryName: stream.GameName,
		Title:        stream.Title,
		StartedAt:    &stream.StartedAt,
	}

	preloads := GetPreloads(ctx)

	for _, preload := range preloads {
		switch preload {
		case "followers":
			result.Followers, err = r.deps.CachedTwitchClient.GetChannelFollowersCountByChannelId(
				ctx,
				dashboardId,
			)
			if err != nil {
				return nil, fmt.Errorf("get channel followers count: %w", err)
			}
		case "subs":
			result.Subs, err = r.deps.CachedTwitchClient.GetChannelSubscribersCountByChannelId(
				ctx,
				dashboardId,
			)
			if err != nil {
				return nil, fmt.Errorf("get channel subs count: %w", err)
			}
		case "title":
		case "categoryName":
		case "categoryId":
			if len(stream.ID) == 0 {
				continue
			}

			channelInformation, err := r.deps.CachedTwitchClient.GetChannelInformationById(
				ctx,
				dashboardId,
			)
			if err != nil {
				return nil, fmt.Errorf("get channel information: %w", err)
			}

			if channelInformation == nil {
				return nil, fmt.Errorf("channel information is nil")
			}

			result.CategoryName = channelInformation.GameName
			result.Title = channelInformation.Title
			result.CategoryID = channelInformation.GameID
		}
	}

	if len(stream.ID) == 0 {
		return result, nil
	}

	parsedMessages, err := r.deps.Redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(stream.ID),
	).Int()
	if err != nil {
		return nil, fmt.Errorf("get stream parsed messages: %w", err)
	}

	result.ChatMessages = parsedMessages

	var (
		usedEmotes     int64
		requestedSongs int64
	)

	if err = r.deps.Gorm.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&usedEmotes).Error; err != nil {
		return nil, fmt.Errorf("get count of used emotes: %w", err)
	}

	if err = r.deps.Gorm.
		WithContext(ctx).
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&requestedSongs).Error; err != nil {
		return nil, fmt.Errorf("get count of requested songs: %w", err)
	}

	result.UsedEmotes = int(usedEmotes)
	result.RequestedSongs = int(requestedSongs)

	return result, nil
}
