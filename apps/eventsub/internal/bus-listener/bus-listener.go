package bus_listener

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/apps/eventsub/internal/manager"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type BusListener struct {
	eventSubClient *manager.Manager
	kickSubManager *kick.SubscriptionManager
	gorm           *gorm.DB
	bus            *buscore.Bus
	logger         *slog.Logger
	channelsRepo   channels.Repository
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
	channelUUID, err := uuid.Parse(msg.ChannelID)
	if err != nil {
		c.logger.Error("error parsing channel ID as UUID", slog.String("channel_id", msg.ChannelID))
		return struct{}{}, fmt.Errorf("parse channel UUID: %w", err)
	}

	return c.subscribeToAllEventsByChannelID(ctx, channelUUID)
}

func (c *BusListener) subscribeToAllEventsByChannelID(
	ctx context.Context,
	channelUUID uuid.UUID,
) (struct{}, error) {
	channelID := channelUUID.String()

	channel, err := c.channelsRepo.GetByID(ctx, channelUUID)
	if err != nil {
		c.logger.Error("error getting channel by ID", slog.String("channel_id", channelID))
		return struct{}{}, err
	}

	if channel.BotID == "" || !channel.IsEnabled {
		c.logger.Warn(
			"channel is not enabled or bot ID is missing",
			slog.String("channel_id", channelID),
		)
		return struct{}{}, nil
	}

	hasActiveSubscription := false

	if channel.KickBotJoined() {
		if channel.KickBotID == nil {
			c.logger.Warn(
				"channel has kick user but no kick bot assigned, skipping kick eventsub subscription",
				slog.String("channel_id", channelID),
				slog.String("kick_user_id", channel.KickUserID.String()),
			)
		} else {
			kickUserIDStr := channel.KickUserID.String()
			if err := c.kickSubManager.SubscribeAll(ctx, *channel.KickUserID); err != nil {
				c.logger.Error(
					"error subscribing to kick events",
					logger.Error(err),
					slog.String("channel_id", channelID),
					slog.String("kick_user_id", kickUserIDStr),
				)
				return struct{}{}, err
			}

			c.logger.Info(
				"subscribed to kick events",
				slog.String("channel_id", channelID),
				slog.String("kick_user_id", kickUserIDStr),
				slog.String("kick_bot_id", channel.KickBotID.String()),
			)
			hasActiveSubscription = true
		}
	}

	if channel.TwitchBotJoined() {
		if channel.TwitchUserID == nil {
			c.logger.Warn(
				"channel has no platform user ID for eventsub subscription",
				slog.String("channel_id", channelID),
			)
			return struct{}{}, nil
		}

		var topics []model.EventsubTopic
		if err := c.gorm.WithContext(ctx).Find(&topics).Error; err != nil {
			c.logger.Error("error getting topics", slog.String("error", err.Error()))
			return struct{}{}, err
		}

		if err := c.eventSubClient.SubscribeToNeededEvents(
			ctx,
			topics,
			*channel.TwitchPlatformID,
			channel.BotID,
		); err != nil {
			return struct{}{}, err
		}

		hasActiveSubscription = true
	}

	if !hasActiveSubscription {
		c.logger.Warn(
			"channel has no active platform bot subscriptions",
			slog.String("channel_id", channelID),
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

	enabledChannels, err := c.channelsRepo.GetMany(ctx, channels.GetManyInput{
		Enabled: lo.ToPtr(true),
	})
	if err != nil {
		c.logger.Error("error getting channels", logger.Error(err))
		return struct{}{}, err
	}

	var wg sync.WaitGroup

	for _, channel := range enabledChannels {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if _, err := c.subscribeToAllEventsByChannelID(ctx, channel.ID); err != nil {
				c.logger.Error("error subscribing to all events", logger.Error(err))
			}
		}()
	}

	wg.Wait()

	c.logger.Info("reinitialized channels for eventsub", slog.Int("count", len(enabledChannels)))

	return struct{}{}, nil
}

func (c *BusListener) unsubscribe(ctx context.Context, userId string) (struct{}, error) {
	channelUUID, err := uuid.Parse(userId)
	if err != nil {
		c.logger.Error("error parsing channel ID for kick unsubscribe", slog.String("channel_id", userId))
		return struct{}{}, fmt.Errorf("parse channel UUID: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, channelUUID)
	if err != nil {
		c.logger.Error("error getting channel for kick unsubscribe", slog.String("channel_id", userId), logger.Error(err))
		return struct{}{}, err
	}

	if channel.TwitchPlatformID != nil {
		if err := c.eventSubClient.UnsubscribeChannel(ctx, *channel.TwitchPlatformID); err != nil {
			c.logger.Error("error unsubscribe twitch channel", logger.Error(err))
			return struct{}{}, err
		}
	}

	if channel.KickUserID != nil {
		kickUserIDStr := channel.KickUserID.String()
		if err := c.kickSubManager.UnsubscribeAll(ctx, *channel.KickUserID); err != nil {
			c.logger.Error(
				"error unsubscribing from kick events",
				logger.Error(err),
				slog.String("channel_id", userId),
				slog.String("kick_user_id", kickUserIDStr),
			)
			return struct{}{}, err
		}

		c.logger.Info(
			"unsubscribed from kick events",
			slog.String("channel_id", userId),
			slog.String("kick_user_id", kickUserIDStr),
		)
	}

	return struct{}{}, nil
}
