package resolvers

import (
	"context"
	"fmt"
	"slices"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/libs/redis_keys"
)

func (r *Resolver) getDashboardStats(ctx context.Context) (*gqlmodel.DashboardStats, error) {
	preloads := GetPreloads(ctx)

	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var stream model.ChannelsStreams
	if err := r.gorm.WithContext(ctx).Where(
		`"userId" = ?`,
		dashboardId,
	).Find(&stream).Error; err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	result := &gqlmodel.DashboardStats{
		Viewers:      &stream.ViewerCount,
		CategoryID:   stream.GameId,
		CategoryName: stream.GameName,
		Title:        stream.Title,
	}

	if slices.Contains(preloads, "followers") {
		result.Followers, err = r.cachedTwitchClient.GetChannelFollowersCountByChannelId(
			ctx,
			dashboardId,
		)
		if err != nil {
			return nil, err
		}
	}

	if slices.Contains(preloads, "title") ||
		slices.Contains(preloads, "categoryName") ||
		slices.Contains(preloads, "categoryId") && stream.ID == "" {
		channelInformation, err := r.cachedTwitchClient.GetChannelInformationById(ctx, dashboardId)
		if err != nil {
			return nil, err
		}

		if channelInformation == nil {
			return nil, fmt.Errorf("channel information is nil")
		}

		result.CategoryName = channelInformation.GameName
		result.Title = channelInformation.Title
		result.CategoryID = channelInformation.GameID
	}
	if stream.ID != "" {
		parsedMessages, _ := r.redis.Get(
			ctx,
			redis_keys.StreamParsedMessages(
				stream.ID,
			),
		).Int()

		result.ChatMessages = parsedMessages
		result.StartedAt = &stream.StartedAt
	}

	var usedEmotes int64
	if err := r.gorm.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&usedEmotes).Error; err != nil {
		return nil, fmt.Errorf("failed to get used emotes: %w", err)
	}

	var requestedSongs int64
	if err := r.gorm.
		WithContext(ctx).
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&requestedSongs).Error; err != nil {
		return nil, fmt.Errorf("failed to get requested songs: %w", err)
	}

	result.UsedEmotes = int(usedEmotes)
	result.RequestedSongs = int(requestedSongs)

	if slices.Contains(preloads, "subs") {
		result.Subs, err = r.cachedTwitchClient.GetChannelSubscribersCountByChannelId(ctx, dashboardId)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
