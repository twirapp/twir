package bus_listener

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/nicklaw5/helix/v2"
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
	gorm           *gorm.DB
	bus            *buscore.Bus
	logger         logger.Logger
	channelsRepo   channels.Repository
	config         config.Config
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Manager      *manager.Manager
	Gorm         *gorm.DB
	Bus          *buscore.Bus
	Logger       logger.Logger
	ChannelsRepo channels.Repository
	Config       config.Config
}

func New(opts Opts) (*BusListener, error) {
	impl := &BusListener{
		eventSubClient: opts.Manager,
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
	channel, err := c.channelsRepo.GetByID(ctx, msg.ChannelID)
	if err != nil {
		c.logger.Error("error getting channel by ID", err, slog.String("channel_id", msg.ChannelID))
		return struct{}{}, err
	}

	if channel.BotID == "" || !channel.IsEnabled {
		c.logger.Warn(
			"channel is not enabled or bot ID is missing",
			slog.String("channel_id", msg.ChannelID),
		)
		return struct{}{}, nil
	}

	var topics []model.EventsubTopic
	if err := c.gorm.WithContext(ctx).Find(&topics).Error; err != nil {
		c.logger.Error("error getting topics", err)
		return struct{}{}, err
	}

	if err := c.eventSubClient.SubscribeToNeededEvents(
		ctx,
		topics,
		msg.ChannelID,
		channel.BotID,
	); err != nil {
		return struct{}{}, err
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
		c.logger.Error("error subscribing to event", err)
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
		c.logger.Error("error creating Twitch app client", err)
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
			c.logger.Error("error getting subscriptions from Twitch", err)
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
					c.logger.Error("error removing subscription", err, slog.String("subscription_id", sub.ID))
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

	var ch []model.Channels
	err = c.gorm.
		WithContext(ctx).
		Select("id").
		Where(`"isEnabled" = true`).Find(&ch).Error
	if err != nil {
		c.logger.Error("error getting channels", err)
		return struct{}{}, err
	}

	var wg sync.WaitGroup

	for _, channel := range ch {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if _, err := c.subscribeToAllEvents(
				ctx,
				eventsub.EventsubSubscribeToAllEventsRequest{
					ChannelID: channel.ID,
				},
			); err != nil {
				c.logger.Error("error subscribing to all events", slog.Any("err", err))
			}
		}()
	}

	wg.Wait()

	c.logger.Info("reinitialized channels for eventsub", slog.Int("count", len(ch)))

	return struct{}{}, nil
}

func (c *BusListener) unsubscribe(ctx context.Context, userId string) (struct{}, error) {
	if err := c.eventSubClient.UnsubscribeChannel(ctx, userId); err != nil {
		c.logger.Error("error unsubscribe channel", err)
		return struct{}{}, err
	}

	return struct{}{}, nil
}
