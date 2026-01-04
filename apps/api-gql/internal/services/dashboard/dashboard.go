package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm                    *gorm.DB
	CachedTwitchClient      *twitchcache.CachedTwitchClient
	KV                      kv.KV
	Config                  config.Config
	Logger                  *slog.Logger
	TwirBus                 *buscore.Bus
	ChannelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	ChannelEmotesUsagesRepo channelsemotesusagesrepository.Repository
	StreamsRepository       streams.Repository
}

func New(opts Opts) *Service {
	return &Service{
		gorm:                    opts.Gorm,
		cachedTwitchClient:      opts.CachedTwitchClient,
		kv:                      opts.KV,
		config:                  opts.Config,
		logger:                  opts.Logger,
		twirBus:                 opts.TwirBus,
		channelsCache:           opts.ChannelsCache,
		channelEmotesUsagesRepo: opts.ChannelEmotesUsagesRepo,
		streamsRepository:       opts.StreamsRepository,
	}
}

type Service struct {
	gorm                    *gorm.DB
	cachedTwitchClient      *twitchcache.CachedTwitchClient
	kv                      kv.KV
	config                  config.Config
	logger                  *slog.Logger
	twirBus                 *buscore.Bus
	channelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	channelEmotesUsagesRepo channelsemotesusagesrepository.Repository
	streamsRepository       streams.Repository
}

func (c *Service) GetDashboardStats(ctx context.Context, channelID string) (
	*entity.DashboardStats,
	error,
) {
	stream, err := c.streamsRepository.GetByChannelID(
		ctx,
		channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("get stream by channel id: %w", err)
	}

	if stream.IsNil() {
		return nil, nil
	}

	result := entity.DashboardStats{
		StreamViewers:      &stream.ViewerCount,
		StreamCategoryID:   stream.GameId,
		StreamCategoryName: stream.GameName,
		StreamTitle:        stream.Title,
		StreamStartedAt:    &stream.StartedAt,
	}

	channelTwitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		channelID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel twitch client: %w", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		followers, err := channelTwitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: channelID,
			},
		)
		if err != nil {
			c.logger.Error("cannot get followers", logger.Error(err))
			return
		}
		if followers.ErrorMessage != "" {
			c.logger.Error("cannot get followers", slog.String("error", followers.ErrorMessage))
			return
		}

		result.Followers = followers.Data.Total
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
		return &result, nil
	}

	parsedMessages, _ := c.kv.Get(
		ctx,
		redis_keys.StreamParsedMessages(stream.ID),
	).Int()

	result.StreamChatMessages = int(parsedMessages)

	var (
		usedEmotes     int64
		requestedSongs int64
	)

	var errgrp errgroup.Group
	errgrp.Go(
		func() error {
			emotesCount, err := c.channelEmotesUsagesRepo.Count(
				ctx,
				channelsemotesusagesrepository.CountInput{
					ChannelID: &channelID,
					TimeAfter: &stream.StartedAt,
				},
			)
			if err != nil {
				return fmt.Errorf("get count of used emotes: %w", err)
			}

			usedEmotes = int64(emotesCount)

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
		return nil, err
	}

	result.UsedEmotes = int(usedEmotes)
	result.RequestedSongs = int(requestedSongs)

	return &result, nil
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

	twitchClient, err := twitch.NewUserClientWithContext(ctx, channelID, c.config, c.twirBus)
	if err != nil {
		return entity.BotStatus{}, err
	}

	result := entity.BotStatus{
		Enabled: dbUser.Channel.IsEnabled,
		IsMod:   false,
	}

	var errgrp errgroup.Group

	errgrp.Go(
		func() error {
			if channelID == dbUser.Channel.BotID {
				result.IsMod = true
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
				result.IsMod = true
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

const (
	BotJoinLeaveActionJoin  = "JOIN"
	BotJoinLeaveActionLeave = "LEAVE"
)

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

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
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
			ctx,
			eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID},
		)
	}

	broadcasterClient, err := twitch.NewUserClientWithContext(
		ctx,
		channelID,
		c.config,
		c.twirBus,
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

	c.channelsCache.Invalidate(ctx, channelID)

	return true, nil
}
