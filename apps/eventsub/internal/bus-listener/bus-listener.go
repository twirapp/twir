package bus_listener

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/apps/eventsub/internal/manager"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type eventSubManager interface {
	UnsubscribeChannel(context.Context, string) error
	SubscribeToNeededEvents(context.Context, []model.EventsubTopic, string, string) error
	SubscribeToEvent(context.Context, string, string, string) error
}

type kickSubscriptionManager interface {
	Subscribe(context.Context, channelplatformentity.ChannelPlatform) error
	Unsubscribe(context.Context, channelplatformentity.ChannelPlatform) error
}

type channelReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error)
}

type BusListener struct {
	eventSubClient eventSubManager
	kickSubManager kickSubscriptionManager
	gorm           *gorm.DB
	bus            *buscore.Bus
	logger         *slog.Logger
	channelsRepo   channels.Repository
	channelService channelReader
	config         config.Config
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Manager        *manager.Manager
	KickSubManager *kick.SubscriptionManager
	Gorm           *gorm.DB
	Bus            *buscore.Bus
	Logger         *slog.Logger
	ChannelsRepo   channels.Repository
	ChannelService *channelservice.ChannelService
	Config         config.Config
}

func New(opts Opts) (*BusListener, error) {
	impl := &BusListener{
		eventSubClient: opts.Manager,
		kickSubManager: opts.KickSubManager,
		gorm:           opts.Gorm,
		bus:            opts.Bus,
		logger:         opts.Logger,
		channelsRepo:   opts.ChannelsRepo,
		channelService: opts.ChannelService,
		config:         opts.Config,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.bus.EventSub.SubscribeToAllEvents.SubscribeGroup(
					"eventsub",
					impl.subscribeToAllEvents,
				); err != nil {
					return err
				}

				if err := impl.bus.EventSub.Subscribe.SubscribeGroup(
					"eventsub",
					impl.subscribe,
				); err != nil {
					return err
				}

				if err := impl.bus.EventSub.InitChannels.SubscribeGroup(
					"eventsub",
					impl.reinitChannels,
				); err != nil {
					return err
				}

				if err := impl.bus.EventSub.Unsubscribe.SubscribeGroup(
					"eventsub",
					impl.unsubscribe,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.bus.EventSub.SubscribeToAllEvents.Unsubscribe()
				impl.bus.EventSub.Subscribe.Unsubscribe()
				impl.bus.EventSub.InitChannels.Unsubscribe()
				return nil
			},
		},
	)

	return impl, nil
}

func (c *BusListener) subscribeToAllEvents(
	ctx context.Context,
	msg eventsub.EventsubSubscribeToAllEventsRequest,
) (struct{}, error) {
	if msg.Platform != "" && !msg.Platform.IsValid() {
		return struct{}{}, fmt.Errorf("invalid platform: %s", msg.Platform)
	}

	platformLabel := msg.Platform.String()
	if platformLabel == "" {
		platformLabel = "all"
	}

	c.logger.Info(
		"received subscribe to all events request",
		slog.String("channel_id", msg.ChannelID),
		slog.String("platform", platformLabel),
	)

	channelUUID, err := uuid.Parse(msg.ChannelID)
	if err != nil {
		c.logger.Error("error parsing channel ID as UUID", slog.String("channel_id", msg.ChannelID))
		return struct{}{}, fmt.Errorf("parse channel UUID: %w", err)
	}

	return c.subscribeToAllEventsByChannelID(ctx, channelUUID, msg.Platform)
}

func (c *BusListener) subscribeToAllEventsByChannelID(
	ctx context.Context,
	channelUUID uuid.UUID,
	platform platformentity.Platform,
) (struct{}, error) {
	channelID := channelUUID.String()
	platformLabel := platform.String()
	if platformLabel == "" {
		platformLabel = "all"
	}

	channel, err := c.channelService.GetChannelByID(ctx, channelUUID)
	if err != nil {
		c.logger.Error("error getting channel by ID", slog.String("channel_id", channelID))
		return struct{}{}, err
	}

	twitchBinding, hasTwitchBinding := channel.Binding(platformentity.PlatformTwitch)
	hasTwitchBinding = hasTwitchBinding &&
		twitchBinding.UserID != uuid.Nil &&
		twitchBinding.PlatformChannelID != ""
	kickBinding, hasKickBinding := channel.Binding(platformentity.PlatformKick)
	hasKickBinding = hasKickBinding &&
		kickBinding.UserID != uuid.Nil &&
		kickBinding.PlatformChannelID != ""

	var twitchBotID string
	if (platform == "" || platform == platformentity.PlatformTwitch) && hasTwitchBinding {
		botConfig, configErr := twitchBinding.ParseTwitchBotConfig()
		if configErr != nil {
			c.logger.Error(
				"cannot parse Twitch bot config",
				logger.Error(configErr),
				slog.String("channel_id", channelID),
			)
		} else {
			twitchBotID = botConfig.BotID
		}
	}

	if platform == platformentity.PlatformTwitch &&
		(!hasTwitchBinding || !twitchBinding.Enabled || twitchBotID == "") {
		c.logger.Warn(
			"channel bot is not enabled for platform",
			slog.String("channel_id", channelID),
			slog.String("platform", platformLabel),
		)
		return struct{}{}, nil
	}
	if platform == platformentity.PlatformKick && (!hasKickBinding || !kickBinding.Enabled) {
		c.logger.Warn(
			"channel bot is not enabled for platform",
			slog.String("channel_id", channelID),
			slog.String("platform", platformLabel),
		)
		return struct{}{}, nil
	}

	// Unsubscribe first (idempotent) to prevent race condition where
	// a separate Unsubscribe bus message arrives after we subscribe.
	if (platform == "" || platform == platformentity.PlatformTwitch) && hasTwitchBinding {
		if err := c.eventSubClient.UnsubscribeChannel(ctx, twitchBinding.PlatformChannelID); err != nil {
			c.logger.Warn("error unsubscribing twitch before resubscribe (continuing)",
				logger.Error(err), slog.String("channel_id", channelID))
		}
	}

	hasActiveSubscription := false

	if (platform == "" || platform == platformentity.PlatformKick) && hasKickBinding && kickBinding.Enabled {
		if kickBinding.BotUserID == nil {
			c.logger.Warn(
				"channel has kick user but no kick bot assigned, skipping kick eventsub subscription",
				slog.String("channel_id", channelID),
				slog.String("platform", platformLabel),
				slog.String("kick_user_id", kickBinding.UserID.String()),
			)
		} else {
			kickUserIDStr := kickBinding.UserID.String()
			if err := c.kickSubManager.Subscribe(ctx, kickBinding); err != nil {
				c.logger.Error(
					"error subscribing to kick events",
					logger.Error(err),
					slog.String("channel_id", channelID),
					slog.String("platform", platformLabel),
					slog.String("kick_user_id", kickUserIDStr),
				)
				return struct{}{}, err
			}

			c.logger.Info(
				"subscribed to kick events",
				slog.String("channel_id", channelID),
				slog.String("platform", platformLabel),
				slog.String("kick_user_id", kickUserIDStr),
				slog.String("kick_bot_user_id", kickBinding.BotUserID.String()),
			)
			hasActiveSubscription = true
		}
	}

	if (platform == "" || platform == platformentity.PlatformTwitch) && hasTwitchBinding && twitchBinding.Enabled {
		if twitchBotID == "" {
			c.logger.Warn(
				"channel bot ID is missing",
				slog.String("channel_id", channelID),
				slog.String("platform", platformLabel),
			)
		} else {
			var topics []model.EventsubTopic
			if err := c.gorm.WithContext(ctx).Find(&topics).Error; err != nil {
				c.logger.Error("error getting topics", slog.String("error", err.Error()), slog.String("platform", platformLabel))
				return struct{}{}, err
			}

			if err := c.eventSubClient.SubscribeToNeededEvents(
				ctx,
				topics,
				twitchBinding.PlatformChannelID,
				twitchBotID,
			); err != nil {
				return struct{}{}, err
			}

			hasActiveSubscription = true
		}
	}

	if !hasActiveSubscription {
		c.logger.Warn(
			"channel has no active platform bot subscriptions",
			slog.String("channel_id", channelID),
			slog.String("platform", platformLabel),
		)
		return struct{}{}, nil
	}

	return struct{}{}, nil
}

func (c *BusListener) subscribe(
	ctx context.Context,
	msg eventsub.EventsubSubscribeRequest,
) (struct{}, error) {
	if err := c.eventSubClient.SubscribeToEvent(
		ctx,
		msg.Topic,
		msg.Version,
		msg.ChannelID,
	); err != nil {
		c.logger.Error("error subscribing to event", logger.Error(err))
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (c *BusListener) reinitChannels(
	ctx context.Context,
	_ struct{},
) (struct{}, error) {
	ctx = context.WithoutCancel(ctx)

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.bus)
	if err != nil {
		c.logger.Error("error creating Twitch app client", logger.Error(err))
		return struct{}{}, err
	}

	var i atomic.Int64
	var cursor string
	for {
		subs, err := twitchClient.GetEventSubSubscriptions(
			&helix.EventSubSubscriptionsParams{
				After: cursor,
			},
		)
		if err != nil {
			c.logger.Error("error getting subscriptions from Twitch", logger.Error(err))
			return struct{}{}, err
		}
		if subs.ErrorMessage != "" {
			c.logger.Error("error in Twitch response", slog.String("error", subs.ErrorMessage))
			return struct{}{}, fmt.Errorf("error getting subscriptions: %s", subs.ErrorMessage)
		}

		var wg sync.WaitGroup

		for _, sub := range subs.Data.EventSubSubscriptions {
			wg.Add(1)

			go func() {
				defer wg.Done()
				resp, err := twitchClient.RemoveEventSubSubscription(sub.ID)
				if err != nil {
					c.logger.Error("error removing subscription", logger.Error(err), slog.String("subscription_id", sub.ID))
					return
				}
				if resp.ErrorMessage != "" {
					c.logger.Error(
						"error in Twitch response while removing subscription",
						slog.String("error", resp.ErrorMessage),
						slog.String("subscription_id", sub.ID),
					)
					return
				}

				i.Add(1)
				c.logger.Info(
					"removed subscription",
					slog.String("subscription_id", sub.ID),
					slog.Int64("removed_count", i.Load()),
				)
			}()
		}

		wg.Wait()

		cursor = subs.Data.Pagination.Cursor
		if cursor == "" {
			break
		}
	}

	reinitializedChannels, err := c.reinitBoundChannels(ctx, func(channelID uuid.UUID) {
		if _, err := c.subscribeToAllEventsByChannelID(ctx, channelID, ""); err != nil {
			c.logger.Error("error subscribing to all events", logger.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("error getting channels", logger.Error(err))
		return struct{}{}, err
	}

	c.logger.Info("reinitialized channels for eventsub", slog.Int("count", reinitializedChannels))

	return struct{}{}, nil
}

func (c *BusListener) reinitBoundChannels(ctx context.Context, reinit func(uuid.UUID)) (int, error) {
	const maxConcurrentReinitializations = 10

	twitchChannels, err := c.channelsRepo.GetAllByBindingPlatform(ctx, platformentity.PlatformTwitch)
	if err != nil {
		return 0, err
	}

	kickChannels, err := c.channelsRepo.GetAllByBindingPlatform(ctx, platformentity.PlatformKick)
	if err != nil {
		return 0, err
	}

	channelIDs := make(map[uuid.UUID]struct{}, len(twitchChannels)+len(kickChannels))
	for _, channel := range twitchChannels {
		channelIDs[channel.ID] = struct{}{}
	}
	for _, channel := range kickChannels {
		channelIDs[channel.ID] = struct{}{}
	}

	semaphore := make(chan struct{}, maxConcurrentReinitializations)
	var wg sync.WaitGroup
	for channelID := range channelIDs {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(channelID uuid.UUID) {
			defer wg.Done()
			defer func() { <-semaphore }()
			reinit(channelID)
		}(channelID)
	}

	wg.Wait()

	return len(channelIDs), nil
}

func (c *BusListener) unsubscribe(ctx context.Context, msg eventsub.EventsubUnsubscribeRequest) (struct{}, error) {
	if msg.Binding != nil {
		return c.unsubscribeSnapshot(ctx, msg)
	}

	channelUUID, err := uuid.Parse(msg.ChannelID)
	if err != nil {
		c.logger.Error("error parsing channel ID for unsubscribe", slog.String("channel_id", msg.ChannelID))
		return struct{}{}, fmt.Errorf("parse channel UUID: %w", err)
	}

	channel, err := c.channelService.GetChannelByID(ctx, channelUUID)
	if err != nil {
		c.logger.Error("error getting channel for unsubscribe", slog.String("channel_id", msg.ChannelID), logger.Error(err))
		return struct{}{}, err
	}

	twitchBinding, hasTwitchBinding := channel.Binding(platformentity.PlatformTwitch)
	if (msg.Platform == "" || msg.Platform == platformentity.PlatformTwitch) &&
		hasTwitchBinding && twitchBinding.PlatformChannelID != "" {
		if err := c.eventSubClient.UnsubscribeChannel(ctx, twitchBinding.PlatformChannelID); err != nil {
			c.logger.Error("error unsubscribe twitch channel", logger.Error(err))
			return struct{}{}, err
		}
	}

	kickBinding, hasKickBinding := channel.Binding(platformentity.PlatformKick)
	if (msg.Platform == "" || msg.Platform == platformentity.PlatformKick) && hasKickBinding {
		kickUserIDStr := kickBinding.UserID.String()
		if err := c.kickSubManager.Unsubscribe(ctx, kickBinding); err != nil {
			c.logger.Error(
				"error unsubscribing from kick events",
				logger.Error(err),
				slog.String("channel_id", msg.ChannelID),
				slog.String("kick_user_id", kickUserIDStr),
			)
			return struct{}{}, err
		}

		c.logger.Info(
			"unsubscribed from kick events",
			slog.String("channel_id", msg.ChannelID),
			slog.String("kick_user_id", kickUserIDStr),
		)
	}

	return struct{}{}, nil
}

func (c *BusListener) unsubscribeSnapshot(
	ctx context.Context,
	msg eventsub.EventsubUnsubscribeRequest,
) (struct{}, error) {
	switch msg.Platform {
	case platformentity.PlatformTwitch:
		if msg.Binding.PlatformChannelID == "" {
			return struct{}{}, fmt.Errorf("missing Twitch platform channel ID for unsubscribe")
		}
		if err := c.eventSubClient.UnsubscribeChannel(ctx, msg.Binding.PlatformChannelID); err != nil {
			return struct{}{}, fmt.Errorf("unsubscribe Twitch channel: %w", err)
		}
	case platformentity.PlatformKick:
		bindingID, err := uuid.Parse(msg.Binding.ID)
		if err != nil {
			return struct{}{}, fmt.Errorf("parse Kick binding ID: %w", err)
		}
		userID, err := uuid.Parse(msg.Binding.UserID)
		if err != nil {
			return struct{}{}, fmt.Errorf("parse Kick binding user ID: %w", err)
		}
		if err := c.kickSubManager.Unsubscribe(ctx, channelplatformentity.ChannelPlatform{
			ID:                bindingID,
			Platform:          platformentity.PlatformKick,
			UserID:            userID,
			PlatformChannelID: msg.Binding.PlatformChannelID,
		}); err != nil {
			return struct{}{}, fmt.Errorf("unsubscribe Kick channel: %w", err)
		}
	}

	return struct{}{}, nil
}
