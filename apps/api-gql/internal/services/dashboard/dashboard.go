package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
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
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	"github.com/twirapp/twir/libs/repositories/streams"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
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
	ChannelsRepo            channelsrepository.Repository
	ChannelEmotesUsagesRepo channelsemotesusagesrepository.Repository
	StreamsRepository       streams.Repository
	UsersRepo               usersrepository.Repository
	KickBotsRepo            kickbotsrepository.Repository
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
		channelsRepo:            opts.ChannelsRepo,
		channelEmotesUsagesRepo: opts.ChannelEmotesUsagesRepo,
		streamsRepository:       opts.StreamsRepository,
		usersRepo:               opts.UsersRepo,
		kickBotsRepo:            opts.KickBotsRepo,
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
	channelsRepo            channelsrepository.Repository
	channelEmotesUsagesRepo channelsemotesusagesrepository.Repository
	streamsRepository       streams.Repository
	usersRepo               usersrepository.Repository
	kickBotsRepo            kickbotsrepository.Repository
}

// getChannelTwitchPlatformID returns the Twitch platform_id (e.g. "12345678") for a channel
// by looking up the user linked as twitch_user_id.
// Returns "" if the channel has no TwitchUserID set.
func (c *Service) getChannelTwitchPlatformID(ctx context.Context, channel channelmodel.Channel) (string, error) {
	if channel.TwitchUserID == nil {
		return "", nil
	}

	user, err := c.usersRepo.GetByID(ctx, channel.TwitchUserID.String())
	if err != nil {
		return "", fmt.Errorf("get user for twitch platform id: %w", err)
	}

	return user.PlatformID, nil
}

func (c *Service) GetDashboardStats(ctx context.Context, channelID string) (
	*entity.DashboardStats,
	error,
) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return nil, fmt.Errorf("channel not found")
	}

	twitchChannelID, err := c.getChannelTwitchPlatformID(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("get channel twitch id: %w", err)
	}

	stream, err := c.streamsRepository.GetByChannelID(
		ctx,
		channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("get stream by channel id: %w", err)
	}

	result := entity.DashboardStats{}

	if twitchChannelID == "" {
		if !stream.IsNil() {
			result.StreamViewers = &stream.ViewerCount
			result.StreamCategoryID = stream.GameId
			result.StreamCategoryName = stream.GameName
			result.StreamTitle = stream.Title
			result.StreamStartedAt = &stream.StartedAt

			parsedMessages, _ := c.kv.Get(
				ctx,
				redis_keys.StreamParsedMessages(stream.ID),
			).Int()
			result.StreamChatMessages = int(parsedMessages)

			var errgrp errgroup.Group
			var usedEmotes int64
			var requestedSongs int64

			errgrp.Go(func() error {
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
			})

			errgrp.Go(func() error {
				if err = c.gorm.
					WithContext(ctx).
					Model(&model.RequestedSong{}).
					Where(`"channelId" = ? AND "createdAt" >= ?`, channelID, stream.StartedAt).
					Count(&requestedSongs).Error; err != nil {
					return fmt.Errorf("get count of requested songs: %w", err)
				}
				return nil
			})

			if err := errgrp.Wait(); err != nil {
				return nil, err
			}

			result.UsedEmotes = int(usedEmotes)
			result.RequestedSongs = int(requestedSongs)
		}

		return &result, nil
	}

	channelTwitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchChannelID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get channel twitch client: %w", err)
	}

	if stream.IsNil() {
		channelInformation, err := channelTwitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
			BroadcasterIDs: []string{twitchChannelID},
		})
		if err != nil {
			return nil, fmt.Errorf("get channel information: %w", err)
		}
		if channelInformation.ErrorMessage != "" {
			return nil, fmt.Errorf("get channel information: %s", channelInformation.ErrorMessage)
		}
		if len(channelInformation.Data.Channels) > 0 {
			c := channelInformation.Data.Channels[0]
			result.StreamCategoryName = c.GameName
			result.StreamTitle = c.Title
			result.StreamCategoryID = c.GameID
		}
	} else {
		result.StreamViewers = &stream.ViewerCount
		result.StreamCategoryID = stream.GameId
		result.StreamCategoryName = stream.GameName
		result.StreamTitle = stream.Title
		result.StreamStartedAt = &stream.StartedAt
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		followers, err := channelTwitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: twitchChannelID,
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
	})

	wg.Go(func() {
		subs, err := c.cachedTwitchClient.GetChannelSubscribersCountByChannelId(
			ctx,
			twitchChannelID,
		)
		if err != nil {
			result.Subs = subs
		}
	})

	if stream.ID == "" {
		wg.Go(func() {
			channelInformation, err := c.cachedTwitchClient.GetChannelInformationById(
				ctx,
				twitchChannelID,
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
		})
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
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return entity.BotStatus{}, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, parsedID)
	if err != nil {
		return entity.BotStatus{}, fmt.Errorf("get channel: %w", err)
	}

	if channel.IsNil() {
		return entity.BotStatus{}, fmt.Errorf("channel not found")
	}

	result := entity.BotStatus{
		Enabled: channel.IsEnabled,
		IsMod:   channel.IsBotMod,
		BotID:   channel.BotID,
	}

	if channel.KickUserID != nil {
		bot, err := c.kickBotsRepo.GetDefault(ctx)
		if err != nil {
			c.logger.Error("cannot get default kick bot", logger.Error(err))
		} else {
			result.BotName = bot.KickUserLogin
			result.BotID = bot.KickUserID.String()
		}
		return result, nil
	}

	twitchChannelID, err := c.getChannelTwitchPlatformID(ctx, channel)
	if err != nil {
		return entity.BotStatus{}, err
	}

	if twitchChannelID == "" {
		return result, nil
	}

	twitchClient, err := twitch.NewUserClientWithContext(ctx, twitchChannelID, c.config, c.twirBus)
	if err != nil {
		return entity.BotStatus{}, err
	}

	var errgrp errgroup.Group

	errgrp.Go(
		func() error {
			if twitchChannelID == channel.BotID {
				result.IsMod = true
				return nil
			}

			mods, err := twitchClient.GetModerators(
				&helix.GetModeratorsParams{
					BroadcasterID: twitchChannelID,
					UserIDs:       []string{channel.BotID},
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
					IDs: []string{channel.BotID},
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

	if _, err := c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{IsBotMod: &result.IsMod}); err != nil {
		c.logger.Error("cannot update channel", logger.Error(err), slog.String("channelId", channel.ID.String()))
	}

	return result, nil
}

const (
	BotJoinLeaveActionJoin  = "JOIN"
	BotJoinLeaveActionLeave = "LEAVE"
)

func (c *Service) BotJoinLeave(ctx context.Context, channelID, action string) (bool, error) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return false, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, parsedID)
	if err != nil {
		return false, err
	}

	if channel.IsNil() {
		return false, fmt.Errorf("channel not found")
	}

	isEnabled := action == BotJoinLeaveActionJoin

	_, err = c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{IsEnabled: &isEnabled})
	if err != nil {
		return false, fmt.Errorf("update channel enabled state: %w", err)
	}

	channel.IsEnabled = isEnabled

	if channel.TwitchUserID == nil {
		// Non-Twitch channel (e.g. Kick) — just invalidate cache and return
		c.channelsCache.Invalidate(ctx, channelID)
		return true, nil
	}

	twitchChannelID, err := c.getChannelTwitchPlatformID(ctx, channel)
	if err != nil {
		return false, err
	}

	if twitchChannelID == "" {
		return false, fmt.Errorf("twitch channel id not found")
	}

	if action == BotJoinLeaveActionJoin {
		channel.IsEnabled = true
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		return false, err
	}

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{IDs: []string{twitchChannelID}},
	)
	if err != nil || twitchUsers.ErrorMessage != "" || len(twitchUsers.Data.Users) == 0 {
		return false, fmt.Errorf("user not found on twitch")
	}

	if channel.IsEnabled {
		c.twirBus.EventSub.SubscribeToAllEvents.Publish(
			ctx,
			eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID},
		)
	}

	broadcasterClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchChannelID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return false, err
	}

	if action == BotJoinLeaveActionJoin {
		unbanResp, err := broadcasterClient.UnbanUser(
			&helix.UnbanUserParams{
				BroadcasterID: twitchChannelID,
				ModeratorID:   twitchChannelID,
				UserID:        channel.BotID,
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
				BroadcasterID: twitchChannelID,
				UserID:        channel.BotID,
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
