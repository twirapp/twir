package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/redis_keys"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm               *gorm.DB
	CachedTwitchClient *twitchcache.CachedTwitchClient
	Redis              *redis.Client
	Config             config.Config
	TokensClient       tokens.TokensClient
	Logger             logger.Logger
	TwirBus            *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		gorm:               opts.Gorm,
		cachedTwitchClient: opts.CachedTwitchClient,
		redis:              opts.Redis,
		config:             opts.Config,
		tokensClient:       opts.TokensClient,
		logger:             opts.Logger,
		twirBus:            opts.TwirBus,
	}
}

type Service struct {
	gorm               *gorm.DB
	cachedTwitchClient *twitchcache.CachedTwitchClient
	redis              *redis.Client
	config             config.Config
	tokensClient       tokens.TokensClient
	logger             logger.Logger
	twirBus            *buscore.Bus
}

func (c *Service) GetDashboardStats(ctx context.Context, channelID string) (
	entity.DashboardStats,
	error,
) {
	var stream model.ChannelsStreams

	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"userId" = ?`,
			channelID,
		).
		Find(&stream).Error; err != nil {
		return entity.DashboardStats{}, fmt.Errorf("get stream: %w", err)
	}

	result := entity.DashboardStats{
		StreamViewers:      &stream.ViewerCount,
		StreamCategoryID:   stream.GameId,
		StreamCategoryName: stream.GameName,
		StreamTitle:        stream.Title,
		StreamStartedAt:    &stream.StartedAt,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		followers, err := c.cachedTwitchClient.GetChannelFollowersCountByChannelId(
			ctx,
			channelID,
		)

		if err != nil {
			result.Followers = followers
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		subs, err := c.cachedTwitchClient.GetChannelSubscribersCountByChannelId(
			ctx,
			channelID,
		)

		if err != nil {
			result.Subs = subs
		}
	}()

	if stream.ID == "" {
		wg.Add(1)

		go func() {
			defer wg.Done()

			channelInformation, err := c.cachedTwitchClient.GetChannelInformationById(
				ctx,
				channelID,
			)

			if err != nil {
				return
			}

			if channelInformation == nil {
				return
			}

			result.StreamCategoryName = channelInformation.GameName
			result.StreamTitle = channelInformation.Title
			result.StreamCategoryID = channelInformation.GameID
		}()
	}

	wg.Wait()

	if len(stream.ID) == 0 {
		return result, nil
	}

	parsedMessages, err := c.redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(stream.ID),
	).Int()
	if err != nil {
		return entity.DashboardStats{}, fmt.Errorf("get stream parsed messages: %w", err)
	}

	result.StreamChatMessages = parsedMessages

	var (
		usedEmotes     int64
		requestedSongs int64
	)

	var errgrp errgroup.Group
	errgrp.Go(
		func() error {
			if err = c.gorm.
				WithContext(ctx).
				Model(&model.ChannelEmoteUsage{}).
				Where(`"channelId" = ? AND "createdAt" >= ?`, channelID, stream.StartedAt).
				Count(&usedEmotes).Error; err != nil {
				return fmt.Errorf("get count of used emotes: %w", err)
			}

			return nil
		},
	)

	errgrp.Go(
		func() error {
			if err = c.gorm.
				WithContext(ctx).
				Model(&model.RequestedSong{}).
				Where(`"channelId" = ? AND "createdAt" >= ?`, channelID, stream.StartedAt).
				Count(&requestedSongs).Error; err != nil {
				return fmt.Errorf("get count of requested songs: %w", err)
			}

			return nil
		},
	)

	if err := errgrp.Wait(); err != nil {
		return entity.DashboardStats{}, err
	}

	result.UsedEmotes = int(usedEmotes)
	result.RequestedSongs = int(requestedSongs)

	return result, nil
}

func (c *Service) GetBotStatus(ctx context.Context, channelID string) (entity.BotStatus, error) {
	dbUser := &model.Users{}
	err := c.gorm.WithContext(ctx).Where("id = ?", channelID).Preload("Channel").First(dbUser).Error
	if err != nil {
		return entity.BotStatus{}, fmt.Errorf("get user: %w", err)
	}

	if dbUser.ID == "" || dbUser.Channel == nil {
		return entity.BotStatus{}, fmt.Errorf("user not found")
	}

	twitchClient, err := twitch.NewUserClientWithContext(ctx, channelID, c.config, c.tokensClient)
	if err != nil {
		return entity.BotStatus{}, err
	}

	result := entity.BotStatus{
		Enabled: dbUser.Channel.IsEnabled,
		IsMod:   true,
	}

	var errgrp errgroup.Group

	errgrp.Go(
		func() error {
			if channelID == dbUser.Channel.BotID {
				return nil
			}

			mods, err := twitchClient.GetModerators(
				&helix.GetModeratorsParams{
					BroadcasterID: channelID,
					UserIDs:       []string{dbUser.Channel.BotID},
				},
			)
			if err != nil {
				return err
			}
			if mods.ErrorMessage != "" {
				return fmt.Errorf("cannot get moderators: %s", mods.ErrorMessage)
			}

			if len(mods.Data.Moderators) > 0 {
				result.IsMod = false
			}

			return nil
		},
	)

	errgrp.Go(
		func() error {
			infoReq, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: []string{dbUser.Channel.BotID},
				},
			)
			if err != nil {
				return err
			}
			if len(infoReq.Data.Users) == 0 {
				return fmt.Errorf("cannot get user info: %s", infoReq.ErrorMessage)
			}

			result.BotID = infoReq.Data.Users[0].ID
			result.BotName = infoReq.Data.Users[0].Login
			return nil
		},
	)

	if err := errgrp.Wait(); err != nil {
		return entity.BotStatus{}, fmt.Errorf("cannot get bot info: %w", err)
	}

	go func() {
		err := c.gorm.Model(&model.Channels{}).Where("id = ?", dbUser.ID).Update(
			`"isBotMod"`,
			result.IsMod,
		).Error
		if err != nil {
			c.logger.Error("cannot update channel", slog.String("channelId", dbUser.ID))
		}
	}()

	return result, nil
}

const BotJoinLeaveActionJoin = "JOIN"
const BotJoinLeaveActionLeave = "LEAVE"

func (c *Service) BotJoinLeave(ctx context.Context, channelID, action string) (bool, error) {
	dbChannel := &model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", channelID).First(dbChannel).Error
	if err != nil {
		return false, err
	}

	if action == BotJoinLeaveActionJoin {
		dbChannel.IsEnabled = true
	} else {
		dbChannel.IsEnabled = false
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.tokensClient)
	if err != nil {
		return false, err
	}

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{IDs: []string{channelID}},
	)
	if err != nil || twitchUsers.ErrorMessage != "" || len(twitchUsers.Data.Users) == 0 {
		return false, fmt.Errorf("user not found on twitch")
	}

	if err := c.gorm.Where(`"id" = ?`, channelID).Select("*").Save(dbChannel).Error; err != nil {
		return false, err
	}

	if dbChannel.IsEnabled {
		c.twirBus.EventSub.SubscribeToAllEvents.Publish(
			eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID},
		)
	}

	broadcasterClient, err := twitch.NewUserClientWithContext(
		ctx,
		channelID,
		c.config,
		c.tokensClient,
	)
	if err != nil {
		return false, err
	}

	if action == BotJoinLeaveActionJoin {
		unbanResp, err := broadcasterClient.UnbanUser(
			&helix.UnbanUserParams{
				BroadcasterID: channelID,
				ModeratorID:   channelID,
				UserID:        dbChannel.BotID,
			},
		)
		if err != nil {
			return false, err
		}

		if unbanResp.ErrorMessage != "" && unbanResp.StatusCode != 400 {
			return false, fmt.Errorf("cannot unban user: %s", unbanResp.ErrorMessage)
		}

		addModResp, err := broadcasterClient.AddChannelModerator(
			&helix.AddChannelModeratorParams{
				BroadcasterID: channelID,
				UserID:        dbChannel.BotID,
			},
		)
		if err != nil {
			return false, err
		}

		if addModResp.ErrorMessage != "" && unbanResp.StatusCode != 400 {
			return false, fmt.Errorf("cannot add channel moderator: %s", addModResp.ErrorMessage)
		}
	}

	return true, nil
}
