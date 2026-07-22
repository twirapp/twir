package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	apiChannelbinding "github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
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
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/repositories/streams"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm                       *gorm.DB
	CachedTwitchClient         *twitchcache.CachedTwitchClient
	AuthService                *auth.Auth
	KV                         kv.KV
	Config                     config.Config
	Logger                     *slog.Logger
	TwirBus                    *buscore.Bus
	ChannelsCache              *generic_cacher.GenericCacher[channelmodel.Channel]
	ChannelPlatformsRepository channelplatforms.Repository
	ChannelService             *channelservice.ChannelService
	ChannelEmotesUsagesRepo    channelsemotesusagesrepository.Repository
	StreamsRepository          streams.Repository
	UsersRepo                  usersrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		gorm:                    opts.Gorm,
		cachedTwitchClient:      opts.CachedTwitchClient,
		authService:             opts.AuthService,
		kv:                      opts.KV,
		config:                  opts.Config,
		logger:                  opts.Logger,
		twirBus:                 opts.TwirBus,
		channelsCache:           opts.ChannelsCache,
		channelPlatformsRepo:    opts.ChannelPlatformsRepository,
		channelService:          opts.ChannelService,
		channelEmotesUsagesRepo: opts.ChannelEmotesUsagesRepo,
		streamsRepository:       opts.StreamsRepository,
		usersRepo:               opts.UsersRepo,
	}
}

type currentPlatformResolver interface {
	GetCurrentPlatform(ctx context.Context) (string, error)
}

type channelLookup interface {
	GetChannelByID(ctx context.Context, id uuid.UUID) (channelmodel.Channel, error)
}

type channelBindingUpdater interface {
	Patch(
		ctx context.Context,
		id uuid.UUID,
		input channelplatforms.PatchInput,
	) (channelplatformsmodel.ChannelPlatform, error)
}

type usersLookup interface {
	GetByID(ctx context.Context, id uuid.UUID) (usersmodel.User, error)
}

type Service struct {
	gorm                    *gorm.DB
	cachedTwitchClient      *twitchcache.CachedTwitchClient
	authService             currentPlatformResolver
	kv                      kv.KV
	config                  config.Config
	logger                  *slog.Logger
	twirBus                 *buscore.Bus
	channelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	channelPlatformsRepo    channelBindingUpdater
	channelService          channelLookup
	channelEmotesUsagesRepo channelsemotesusagesrepository.Repository
	streamsRepository       streams.Repository
	usersRepo               usersLookup
}

func (c *Service) resolveAnalyticsIdentity(ctx context.Context, channel channelmodel.Channel) (string, string) {
	currentPlatform, err := c.authService.GetCurrentPlatform(ctx)
	if err == nil {
		platform := platformentity.Platform(currentPlatform)
		if platform == platformentity.PlatformKick || platform == platformentity.PlatformTwitch {
			if binding, found := apiChannelbinding.Find(channel, platform); found && binding.PlatformChannelID != "" {
				return currentPlatform, binding.PlatformChannelID
			}
		}
	}

	for _, platform := range []platformentity.Platform{
		platformentity.PlatformTwitch,
		platformentity.PlatformKick,
	} {
		if binding, found := apiChannelbinding.Find(channel, platform); found && binding.PlatformChannelID != "" {
			return platform.String(), binding.PlatformChannelID
		}
	}

	return "", ""
}

func (c *Service) GetDashboardStats(ctx context.Context, channelID string) (
	*entity.DashboardStats,
	error,
) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return nil, fmt.Errorf("channel not found")
	}

	stream, err := c.streamsRepository.GetByChannelID(
		ctx,
		parsedID,
		platformentity.PlatformTwitch,
	)
	if err != nil {
		return nil, fmt.Errorf("get stream by channel id: %w", err)
	}

	result := entity.DashboardStats{}
	analyticsPlatform, analyticsPlatformChannelID := c.resolveAnalyticsIdentity(ctx, channel)

	twitchBinding, hasTwitchBinding := apiChannelbinding.Find(channel, platformentity.PlatformTwitch)
	if !hasTwitchBinding || twitchBinding.PlatformChannelID == "" {
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
				if analyticsPlatform == "" || analyticsPlatformChannelID == "" {
					usedEmotes = 0
					return nil
				}

				emotesCount, err := c.channelEmotesUsagesRepo.Count(
					ctx,
					channelsemotesusagesrepository.CountInput{
						Platform:          &analyticsPlatform,
						PlatformChannelID: &analyticsPlatformChannelID,
						TimeAfter:         &stream.StartedAt,
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

	twitchPlatformID := twitchBinding.PlatformChannelID
	if !stream.IsNil() {
		result.StreamViewers = &stream.ViewerCount
		result.StreamCategoryID = stream.GameId
		result.StreamCategoryName = stream.GameName
		result.StreamTitle = stream.Title
		result.StreamStartedAt = &stream.StartedAt
	}

	channelTwitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchBinding.UserID,
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
			twitchBinding.UserID,
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
			if analyticsPlatform == "" || analyticsPlatformChannelID == "" {
				usedEmotes = 0
				return nil
			}

			emotesCount, err := c.channelEmotesUsagesRepo.Count(
				ctx,
				channelsemotesusagesrepository.CountInput{
					Platform:          &analyticsPlatform,
					PlatformChannelID: &analyticsPlatformChannelID,
					TimeAfter:         &stream.StartedAt,
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

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	if channel.IsNil() {
		return nil, fmt.Errorf("channel not found")
	}

	statuses := make([]entity.BotStatus, 0, 2)

	twitchBinding, twitchBotConfig, hasTwitchBinding, err := apiChannelbinding.FindTwitch(channel)
	if err != nil {
		return nil, fmt.Errorf("parse Twitch channel bot config: %w", err)
	}
	if hasTwitchBinding {
		status, err := c.getTwitchBotStatus(ctx, channel, twitchBinding, twitchBotConfig)
		if err != nil {
			c.logger.Error("cannot get twitch bot status", logger.Error(err), slog.String("channelId", channel.ID.String()))
			status = c.getBasicTwitchBotStatus(ctx, channel, twitchBinding, twitchBotConfig)
		}
		statuses = append(statuses, status)
	}

	if kickBinding, hasKickBinding := apiChannelbinding.Find(channel, platformentity.PlatformKick); hasKickBinding {
		statuses = append(statuses, c.getKickBotStatus(ctx, channel, kickBinding))
	}

	if len(statuses) == 0 {
		anyBindingEnabled := false
		for _, binding := range channel.Bindings {
			anyBindingEnabled = anyBindingEnabled || binding.Enabled
		}

		statuses = append(statuses, entity.BotStatus{
			DashboardID: channel.ID.String(),
			Enabled:     anyBindingEnabled,
		})
	}

	return statuses, nil
}

func (c *Service) getKickBotStatus(
	ctx context.Context,
	channel channelmodel.Channel,
	binding channelplatformsmodel.ChannelPlatform,
) entity.BotStatus {
	result := entity.BotStatus{
		DashboardID: channel.ID.String(),
		Platform:    platformentity.PlatformKick.String(),
		ChannelName: c.getChannelName(ctx, &binding.UserID),
		Enabled:     binding.Enabled,
		IsMod:       true,
	}

	if binding.BotUserID != nil {
		result.BotID = binding.BotUserID.String()
		result.BotName = c.getChannelName(ctx, binding.BotUserID)
	}

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

func (c *Service) getBasicTwitchBotStatus(
	ctx context.Context,
	channel channelmodel.Channel,
	binding channelplatformsmodel.ChannelPlatform,
	botConfig apiChannelbinding.TwitchBotConfig,
) entity.BotStatus {
	return entity.BotStatus{
		DashboardID: channel.ID.String(),
		Platform:    platformentity.PlatformTwitch.String(),
		ChannelName: c.getChannelName(ctx, &binding.UserID),
		Enabled:     binding.Enabled,
		IsMod:       botConfig.IsBotMod,
		BotID:       botConfig.BotID,
		BotName:     "TwirBot",
	}
}

func (c *Service) getTwitchBotStatus(
	ctx context.Context,
	channel channelmodel.Channel,
	binding channelplatformsmodel.ChannelPlatform,
	botConfig apiChannelbinding.TwitchBotConfig,
) (entity.BotStatus, error) {
	result := c.getBasicTwitchBotStatus(ctx, channel, binding, botConfig)

	if binding.PlatformChannelID == "" || botConfig.BotID == "" {
		return result, nil
	}

	twitchPlatformID := binding.PlatformChannelID

	twitchClient, err := twitch.NewUserClientWithContext(ctx, binding.UserID, c.config, c.twirBus)
	if err != nil {
		return entity.BotStatus{}, err
	}

	var errgrp errgroup.Group

	errgrp.Go(
		func() error {
			if twitchPlatformID == botConfig.BotID {
				result.IsMod = true
				return nil
			}

			mods, err := twitchClient.GetModerators(
				&helix.GetModeratorsParams{
					BroadcasterID: twitchPlatformID,
					UserIDs:       []string{botConfig.BotID},
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
					IDs: []string{botConfig.BotID},
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

	botConfigPatch, err := json.Marshal(struct {
		IsBotMod bool `json:"is_bot_mod"`
	}{IsBotMod: result.IsMod})
	if err != nil {
		return entity.BotStatus{}, fmt.Errorf("marshal Twitch bot config patch: %w", err)
	}
	if _, err := c.channelPlatformsRepo.Patch(
		ctx,
		binding.ID,
		channelplatforms.PatchInput{BotConfigPatch: botConfigPatch},
	); err != nil {
		c.logger.Error("cannot update Twitch binding", logger.Error(err), slog.String("channelId", channel.ID.String()))
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

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return false, err
	}

	if channel.IsNil() {
		return false, fmt.Errorf("channel not found")
	}

	targetPlatform := platform
	if targetPlatform == "" {
		if _, found := apiChannelbinding.Find(channel, platformentity.PlatformTwitch); found {
			targetPlatform = "twitch"
		} else if _, found := apiChannelbinding.Find(channel, platformentity.PlatformKick); found {
			targetPlatform = "kick"
		} else {
			return false, fmt.Errorf("channel has no connected platform")
		}
	}

	isEnabled := action == BotJoinLeaveActionJoin

	switch targetPlatform {
	case "kick":
		binding, found := apiChannelbinding.Find(channel, platformentity.PlatformKick)
		if !found || binding.PlatformChannelID == "" {
			return false, fmt.Errorf("kick channel id not found")
		}

		_, err = c.channelPlatformsRepo.Patch(
			ctx,
			binding.ID,
			channelplatforms.PatchInput{Enabled: &isEnabled},
		)
		if err != nil {
			return false, fmt.Errorf("update Kick binding enabled state: %w", err)
		}

		if isEnabled {
			c.twirBus.EventSub.SubscribeToAllEvents.Publish(
				ctx,
				eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID, Platform: platformentity.PlatformKick},
			)
		} else {
			c.twirBus.EventSub.Unsubscribe.Publish(
				ctx,
				eventsub.EventsubUnsubscribeRequest{ChannelID: channelID, Platform: platformentity.PlatformKick},
			)
		}

		c.channelsCache.Invalidate(ctx, channelID)
		return true, nil
	case "twitch":
		binding, botConfig, found, bindingErr := apiChannelbinding.FindTwitch(channel)
		if bindingErr != nil {
			return false, fmt.Errorf("parse Twitch channel bot config: %w", bindingErr)
		}
		if !found || binding.PlatformChannelID == "" {
			return false, fmt.Errorf("twitch channel id not found")
		}
		if botConfig.BotID == "" {
			return false, fmt.Errorf("twitch bot id not found")
		}

		if _, err = c.channelPlatformsRepo.Patch(
			ctx,
			binding.ID,
			channelplatforms.PatchInput{Enabled: &isEnabled},
		); err != nil {
			return false, fmt.Errorf("update Twitch binding enabled state: %w", err)
		}

		twitchPlatformID := binding.PlatformChannelID
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

		if isEnabled {
			c.twirBus.EventSub.SubscribeToAllEvents.Publish(
				ctx,
				eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID, Platform: platformentity.PlatformTwitch},
			)
		} else {
			c.twirBus.EventSub.Unsubscribe.Publish(
				ctx,
				eventsub.EventsubUnsubscribeRequest{ChannelID: channelID, Platform: platformentity.PlatformTwitch},
			)
		}

		broadcasterClient, err := twitch.NewUserClientWithContext(
			ctx,
			binding.UserID,
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
					UserID:        botConfig.BotID,
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
					UserID:        botConfig.BotID,
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
	default:
		return false, fmt.Errorf("unsupported platform: %s", targetPlatform)
	}
}
