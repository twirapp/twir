package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/dashboard"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Dashboard struct {
	*impl_deps.Deps
}

func (c *Dashboard) GetDashboardStats(
	ctx context.Context,
	_ *emptypb.Empty,
) (*dashboard.DashboardStatsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var stream model.ChannelsStreams
	if err := c.Db.WithContext(ctx).Where(
		`"userId" = ?`,
		dashboardId,
	).Find(&stream).Error; err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	twitchClient, err := twitch.NewUserClient(dashboardId, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitch client: %w", err)
	}

	startedAt := fmt.Sprint(stream.StartedAt.UnixMilli())
	var channelCategoryId string
	var channelCategoryName string
	var channelTitle string
	var followersCount uint32
	var subs uint32

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		followsReq, err := twitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: dashboardId,
			},
		)
		if err != nil {
			c.Logger.Error(
				"cannot get followers",
				slog.String("channelId", dashboardId),
				slog.Any("err", err),
			)
		} else if followsReq.ErrorMessage != "" {
			c.Logger.Error(
				"cannot get followers",
				slog.String("channelId", dashboardId),
				slog.Any("err", followsReq.ErrorMessage),
			)
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
			c.Logger.Error(
				"cannot get channel information",
				slog.String("channelId", dashboardId),
				slog.Any("err", err),
			)
		} else if infoReq.ErrorMessage != "" {
			c.Logger.Error(
				"cannot get channel information",
				slog.String("channelId", dashboardId),
				slog.Any("err", infoReq.ErrorMessage),
			)
		} else if len(infoReq.Data.Channels) > 0 {
			channelCategoryId = infoReq.Data.Channels[0].GameID
			channelTitle = infoReq.Data.Channels[0].Title
			channelCategoryName = infoReq.Data.Channels[0].GameName
		}
	}()

	go func() {
		defer wg.Done()
		subsReq, err := twitchClient.GetSubscriptions(
			&helix.SubscriptionsParams{
				BroadcasterID: dashboardId,
			},
		)
		if err != nil {
			c.Logger.Error(
				"cannot get subscriptions",
				slog.String("channelId", dashboardId),
				slog.Any("err", err),
			)
		} else if subsReq.ErrorMessage != "" {
			c.Logger.Error(
				"cannot get subscriptions",
				slog.String("channelId", dashboardId),
				slog.Any("err", subsReq.ErrorMessage),
			)
		} else {
			subs = uint32(subsReq.Data.Total)
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

	return &dashboard.DashboardStatsResponse{
		CategoryId:     channelCategoryId,
		CategoryName:   channelCategoryName,
		Viewers:        uint32(stream.ViewerCount),
		StartedAt:      lo.If[*string](stream.ID == "", nil).Else(&startedAt),
		Title:          channelTitle,
		ChatMessages:   uint32(stream.ParsedMessages),
		Followers:      followersCount,
		UsedEmotes:     uint32(usedEmotes),
		RequestedSongs: 0,
		Subs:           subs,
	}, nil
}

func (c *Dashboard) convertType(t model.ChannelEventListItemType) dashboard.EventType {
	switch t {
	case model.ChannelEventListItemTypeDonation:
		return dashboard.EventType_DONATION
	case model.ChannelEventListItemTypeFollow:
		return dashboard.EventType_FOLLOW
	case model.ChannelEventListItemTypeRaided:
		return dashboard.EventType_RAIDED
	case model.ChannelEventListItemTypeSubscribe:
		return dashboard.EventType_SUBSCRIBE
	case model.ChannelEventListItemTypeReSubscribe:
		return dashboard.EventType_RESUBSCRIBE
	case model.ChannelEventListItemTypeSubGift:
		return dashboard.EventType_SUBGIFT
	case model.ChannelEventListItemTypeFirstUserMessage:
		return dashboard.EventType_FIRST_USER_MESSAGE
	case model.ChannelEventListItemTypeChatClear:
		return dashboard.EventType_CHAT_CLEAR
	case model.ChannelEventListItemTypeRedemptionCreated:
		return dashboard.EventType_REDEMPTION_CREATED
	case model.ChannelEventListItemTypeChannelBan:
		return dashboard.EventType_CHANNEL_BAN
	case model.ChannelEventListItemTypeChannelUnbanRequestCreate:
		return dashboard.EventType_CHANNEL_UNBAN_REQUEST_CREATE
	case model.ChannelEventListItemTypeChannelUnbanRequestResolve:
		return dashboard.EventType_CHANNEL_UNBAN_REQUEST_RESOLVE
	default:
		return 0
	}
}

func (c *Dashboard) GetDashboardEventsList(ctx context.Context, _ *emptypb.Empty) (
	*dashboard.DashboardEventsList,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	var entities []model.ChannelsEventsListItem
	if err := c.Db.Where(
		"channel_id = ?",
		dashboardId,
	).Order("created_at desc").Limit(500).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("cannot get events: %w", err)
	}

	events := make([]*dashboard.DashboardEventsList_Event, len(entities))

	for i, entity := range entities {
		t := c.convertType(entity.Type)

		events[i] = &dashboard.DashboardEventsList_Event{
			UserId: entity.UserID,
			Type:   t,
			Data: &dashboard.EventData{
				DonationAmount:                  entity.Data.DonationAmount,
				DonationCurrency:                entity.Data.DonationCurrency,
				DonationMessage:                 entity.Data.DonationMessage,
				DonationUsername:                entity.Data.DonationUsername,
				RaidedViewersCount:              entity.Data.RaidedViewersCount,
				RaidedFromUserName:              entity.Data.RaidedFromUserName,
				RaidedFromDisplayName:           entity.Data.RaidedFromDisplayName,
				FollowUserName:                  entity.Data.FollowUserName,
				FollowUserDisplayName:           entity.Data.FollowUserDisplayName,
				RedemptionTitle:                 entity.Data.RedemptionTitle,
				RedemptionInput:                 entity.Data.RedemptionInput,
				RedemptionUserName:              entity.Data.RedemptionUserName,
				RedemptionUserDisplayName:       entity.Data.RedemptionUserDisplayName,
				RedemptionCost:                  entity.Data.RedemptionCost,
				SubLevel:                        entity.Data.SubLevel,
				SubUserName:                     entity.Data.SubUserName,
				SubUserDisplayName:              entity.Data.SubUserDisplayName,
				ReSubLevel:                      entity.Data.ReSubLevel,
				ReSubUserName:                   entity.Data.ReSubUserName,
				ReSubUserDisplayName:            entity.Data.ReSubUserDisplayName,
				ReSubMonths:                     entity.Data.ReSubMonths,
				ReSubStreak:                     entity.Data.ReSubStreak,
				SubGiftLevel:                    entity.Data.SubGiftLevel,
				SubGiftUserName:                 entity.Data.SubGiftUserName,
				SubGiftUserDisplayName:          entity.Data.SubGiftUserDisplayName,
				SubGiftTargetUserName:           entity.Data.SubGiftTargetUserName,
				SubGiftTargetUserDisplayName:    entity.Data.SubGiftTargetUserDisplayName,
				FirstUserMessageUserName:        entity.Data.FirstUserMessageUserName,
				FirstUserMessageUserDisplayName: entity.Data.FirstUserMessageUserDisplayName,
				FirstUserMessageMessage:         entity.Data.FirstUserMessageMessage,
				BanReason:                       entity.Data.BanReason,
				BanEndsInMinutes:                entity.Data.BanEndsInMinutes,
				BannedUserName:                  entity.Data.BannedUserName,
				BannedUserLogin:                 entity.Data.BannedUserLogin,
				ModeratorName:                   entity.Data.ModeratorName,
				ModeratorDisplayName:            entity.Data.ModeratorDisplayName,
				Message:                         entity.Data.Message,
				UserLogin:                       entity.Data.UserLogin,
				UserName:                        entity.Data.UserDisplayName,
			},
			CreatedAt: fmt.Sprint(entity.CreatedAt.UnixMilli()),
		}
	}

	return &dashboard.DashboardEventsList{
		Events: events,
	}, nil
}
