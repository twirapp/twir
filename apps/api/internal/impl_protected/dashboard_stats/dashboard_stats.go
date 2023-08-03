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
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
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

	startedAt := fmt.Sprint(stream.StartedAt.UnixMilli())
	var channelCategoryId string
	var channelCategoryName string
	var channelTitle string
	var followersCount uint32

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		followsReq, err := twitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: dashboardId,
			},
		)
		if err != nil {
			zap.S().Error(err)
		} else if followsReq.ErrorMessage != "" {
			zap.S().Error(followsReq.ErrorMessage)
		} else {
			followersCount = uint32(followsReq.Data.Total)
		}
	}()

	go func() {
		defer wg.Done()
		infoReq, err := twitchClient.GetChannelInformation(
			&helix.GetChannelInformationParams{
				BroadcasterIDs: []string{dashboardId},
			},
		)
		if err != nil {
			zap.S().Error(err)
		} else if infoReq.ErrorMessage != "" {
			zap.S().Error(infoReq.ErrorMessage)
		} else if len(infoReq.Data.Channels) > 0 {
			channelCategoryId = infoReq.Data.Channels[0].GameID
			channelTitle = infoReq.Data.Channels[0].Title
			channelCategoryName = infoReq.Data.Channels[0].GameName
		}
	}()

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

	wg.Wait()

	return &dashboard_stats.DashboardStatsResponse{
		CategoryId:     channelCategoryId,
		CategoryName:   channelCategoryName,
		Viewers:        uint32(stream.ViewerCount),
		StartedAt:      lo.If[*string](stream.ID == "", nil).Else(&startedAt),
		Title:          channelTitle,
		ChatMessages:   uint32(stream.ParsedMessages),
		Followers:      followersCount,
		UsedEmotes:     uint32(usedEmotes),
		RequestedSongs: 0,
	}, nil
}
