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
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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

	stream, err := c.streamsRepository.GetByChannelID(
		ctx,
		channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("get stream by channel id: %w", err)
	}

	result := entity.DashboardStats{}

	if !channel.TwitchConnected() {
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

	twitchPlatformID := *channel.TwitchPlatformID
	if !stream.IsNil() {
		result.StreamViewers = &stream.ViewerCount
		result.StreamCategoryID = stream.GameId
		result.StreamCategoryName = stream.GameName
		result.StreamTitle = stream.Title
		result.StreamStartedAt = &stream.StartedAt
	}

	channelTwitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		*channel.TwitchUserID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		appClient, appClientErr := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
		if appClientErr != nil {
			c.logger.Error("cannot get fallback twitch app client", logger.Error(appClientErr))
			return &result, nil
		}

		if stream.IsNil() {
			channelInformation, infoErr := appClient.GetChannelInformation(&helix.GetChannelInformationParams{
				BroadcasterIDs: []string{twitchPlatformID},
			})
			if infoErr != nil {
				c.logger.Error("cannot get channel information with app client", logger.Error(infoErr))
				return &result, nil
			}
			if channelInformation.ErrorMessage != "" {
				c.logger.Error("cannot get channel information with app client", slog.String("error", channelInformation.ErrorMessage))
				return &result, nil
			}
			if len(channelInformation.Data.Channels) > 0 {
				c := channelInformation.Data.Channels[0]
				result.StreamCategoryName = c.GameName
				result.StreamTitle = c.Title
				result.StreamCategoryID = c.GameID
			}
		}

		return &result, nil
	}

	if stream.IsNil() {
		channelInformation, err := channelTwitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
			BroadcasterIDs: []string{twitchPlatformID},
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
				BroadcasterID: twitchPlatformID,
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
			*channel.TwitchUserID,
			twitchPlatformID,
		)
		if err != nil {
			result.Subs = subs
		}
	})

	if stream.ID == "" {
		wg.Go(func() {
			channelInformation, err := c.cachedTwitchClient.GetChannelInformationById(
				ctx,
				twitchPlatformID,
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
	statuses, err := c.GetBotStatuses(ctx, channelID)
	if err != nil {
		return entity.BotStatus{}, err
	}

	if len(statuses) == 0 {
		return entity.BotStatus{}, nil
	}

	return statuses[0], nil
}

func (c *Service) GetBotStatuses(ctx context.Context, channelID string) ([]entity.BotStatus, error) {
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

	statuses := make([]entity.BotStatus, 0, 2)

	if channel.TwitchConnected() {
		status, err := c.getTwitchBotStatus(ctx, channel)
		if err != nil {
			c.logger.Error("cannot get twitch bot status", logger.Error(err), slog.String("channelId", channel.ID.String()))
			status = c.getBasicTwitchBotStatus(ctx, channel)
		}
		statuses = append(statuses, status)
	}

	if channel.KickConnected() {
		statuses = append(statuses, c.getKickBotStatus(ctx, channel))
	}

	if len(statuses) == 0 {
		statuses = append(statuses, entity.BotStatus{
			DashboardID: channel.ID.String(),
			Enabled:     channel.AnyBotJoined(),
			IsMod:       channel.IsBotMod,
			BotID:       channel.BotID,
		})
	}

	return statuses, nil
}

func (c *Service) getKickBotStatus(ctx context.Context, channel channelmodel.Channel) entity.BotStatus {
	result := entity.BotStatus{
		DashboardID: channel.ID.String(),
		Platform:    "kick",
		ChannelName: c.getChannelName(ctx, channel.KickUserID),
		Enabled:     channel.KickBotJoined(),
		IsMod:       true,
		BotID:       channel.BotID,
	}

	bot, err := c.kickBotsRepo.GetDefault(ctx)
	if err != nil {
		c.logger.Error("cannot get default kick bot", logger.Error(err))
		return result
	}

	if channel.KickBotID == nil {
		if _, updateErr := c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{KickBotID: &bot.ID}); updateErr != nil {
			c.logger.Error("cannot repair kick bot assignment", logger.Error(updateErr), slog.String("channelId", channel.ID.String()))
		} else {
			channel.KickBotID = &bot.ID
		}
	}

	result.BotName = bot.KickUserLogin
	result.BotID = bot.KickUserID.String()
	return result
}

func (c *Service) getChannelName(ctx context.Context, userID *uuid.UUID) string {
	if userID == nil {
		return ""
	}

	user, err := c.usersRepo.GetByID(ctx, *userID)
	if err != nil {
		c.logger.Error("cannot get channel user for bot status", logger.Error(err), slog.String("userId", userID.String()))
		return ""
	}

	if user.Login != "" {
		return user.Login
	}

	return user.DisplayName
}

func (c *Service) getBasicTwitchBotStatus(ctx context.Context, channel channelmodel.Channel) entity.BotStatus {
	return entity.BotStatus{
		DashboardID: channel.ID.String(),
		Platform:    "twitch",
		ChannelName: c.getChannelName(ctx, channel.TwitchUserID),
		Enabled:     channel.TwitchBotJoined(),
		IsMod:       channel.IsBotMod,
		BotID:       channel.BotID,
		BotName:     "TwirBot",
	}
}

func (c *Service) getTwitchBotStatus(ctx context.Context, channel channelmodel.Channel) (entity.BotStatus, error) {
	result := c.getBasicTwitchBotStatus(ctx, channel)

	twitchPlatformID := *channel.TwitchPlatformID

	twitchClient, err := twitch.NewUserClientWithContext(ctx, *channel.TwitchUserID, c.config, c.twirBus)
	if err != nil {
		return entity.BotStatus{}, err
	}

	var errgrp errgroup.Group

	errgrp.Go(
		func() error {
			if twitchPlatformID == channel.BotID {
				result.IsMod = true
				return nil
			}

			mods, err := twitchClient.GetModerators(
				&helix.GetModeratorsParams{
					BroadcasterID: twitchPlatformID,
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

func (c *Service) BotJoinLeave(ctx context.Context, channelID, action, platform string) (bool, error) {
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

	targetPlatform := platform
	if targetPlatform == "" {
		switch {
		case channel.TwitchConnected():
			targetPlatform = "twitch"
		case channel.KickConnected():
			targetPlatform = "kick"
		default:
			return false, fmt.Errorf("channel has no connected platform")
		}
	}

	isEnabled := action == BotJoinLeaveActionJoin

	switch targetPlatform {
	case "kick":
		if !channel.KickConnected() {
			return false, fmt.Errorf("kick channel id not found")
		}

		overallEnabled := channel.TwitchBotJoined() || isEnabled
		_, err = c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{
			IsEnabled:      &overallEnabled,
			KickBotEnabled: &isEnabled,
		})
		if err != nil {
			return false, fmt.Errorf("update kick channel enabled state: %w", err)
		}

		if isEnabled {
			c.twirBus.EventSub.SubscribeToAllEvents.Publish(
				ctx,
				eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID, Platform: platformentity.PlatformKick},
			)
		}

		c.channelsCache.Invalidate(ctx, channelID)
		return true, nil
	case "twitch":
		if !channel.TwitchConnected() {
			return false, fmt.Errorf("twitch channel id not found")
		}
	default:
		return false, fmt.Errorf("unsupported platform: %s", targetPlatform)
	}

	overallEnabled := channel.KickBotJoined() || isEnabled
	_, err = c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{
		IsEnabled:        &overallEnabled,
		TwitchBotEnabled: &isEnabled,
	})
	if err != nil {
		return false, fmt.Errorf("update twitch channel enabled state: %w", err)
	}

	channel.IsEnabled = overallEnabled
	channel.TwitchBotEnabled = isEnabled

	twitchPlatformID := *channel.TwitchPlatformID

	if action == BotJoinLeaveActionJoin {
		channel.IsEnabled = true
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		return false, err
	}

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{IDs: []string{twitchPlatformID}},
	)
	if err != nil || twitchUsers.ErrorMessage != "" || len(twitchUsers.Data.Users) == 0 {
		return false, fmt.Errorf("user not found on twitch")
	}

	if channel.TwitchBotJoined() {
		c.twirBus.EventSub.SubscribeToAllEvents.Publish(
			ctx,
			eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID, Platform: platformentity.PlatformTwitch},
		)
	}

	broadcasterClient, err := twitch.NewUserClientWithContext(
		ctx,
		*channel.TwitchUserID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return false, err
	}

	if action == BotJoinLeaveActionJoin {
		unbanResp, err := broadcasterClient.UnbanUser(
			&helix.UnbanUserParams{
				BroadcasterID: twitchPlatformID,
				ModeratorID:   twitchPlatformID,
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
				BroadcasterID: twitchPlatformID,
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
