package dashboard_stats

import (
	"context"
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/dashboard_stats"
	"github.com/satont/twir/libs/twitch"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DashboardStats struct {
	*impl_deps.Deps
}

func (c *DashboardStats) GetDashboardStats(
	ctx context.Context,
	_ *emptypb.Empty,
) (*dashboard_stats.DashboardStatsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var stream model.ChannelsStreams
	if err := c.Db.WithContext(ctx).Where(`"userId" = ?`, dashboardId).Find(&stream).Error; err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	twitchClient, err := twitch.NewUserClient(dashboardId, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitch client: %w", err)
	}

	followsReq, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: dashboardId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel follows: %w", err)
	}
	if followsReq.ErrorMessage != "" {
		return nil, fmt.Errorf("failed to get channel follows: %s", followsReq.ErrorMessage)
	}

	var usedEmotes int64
	if err := c.Db.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&usedEmotes).Error; err != nil {
		return nil, fmt.Errorf("failed to get used emotes: %w", err)
	}

	var requestedSongs int64
	if err := c.Db.
		WithContext(ctx).
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "createdAt" >= ?`, dashboardId, stream.StartedAt).
		Count(&requestedSongs).Error; err != nil {
		return nil, fmt.Errorf("failed to get requested songs: %w", err)
	}

	startedAt := fmt.Sprint(stream.StartedAt.UnixMilli())

	return &dashboard_stats.DashboardStatsResponse{
		CategoryId:     stream.GameId,
		CategoryName:   stream.GameName,
		Viewers:        uint32(stream.ViewerCount),
		StartedAt:      lo.If[*string](stream.ID == "", nil).Else(&startedAt),
		Title:          stream.Title,
		ChatMessages:   uint32(stream.ParsedMessages),
		Followers:      uint32(followsReq.Data.Total),
		UsedEmotes:     uint32(usedEmotes),
		RequestedSongs: 0,
	}, nil
}
