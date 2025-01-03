package bus_listener

import (
	"context"

	"github.com/satont/twir/apps/eventsub/internal/manager"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type BusListener struct {
	eventSubClient *manager.Manager
	gorm           *gorm.DB
	bus            *buscore.Bus
	logger         logger.Logger
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Manager *manager.Manager
	Gorm    *gorm.DB
	Bus     *buscore.Bus
	Logger  logger.Logger
}

func New(opts Opts) (*BusListener, error) {
	impl := &BusListener{
		eventSubClient: opts.Manager,
		gorm:           opts.Gorm,
		bus:            opts.Bus,
		logger:         opts.Logger,
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
) struct{} {
	channel := model.Channels{}
	err := c.gorm.
		WithContext(ctx).
		Where(
			`"id" = ?`,
			msg.ChannelID,
		).First(&channel).Error
	if err != nil {
		c.logger.Error("error getting channel", err)
		return struct{}{}
	}

	var topics []model.EventsubTopic
	if err := c.gorm.WithContext(ctx).Find(&topics).Error; err != nil {
		c.logger.Error("error getting topics", err)
		return struct{}{}
	}

	if err := c.eventSubClient.SubscribeToNeededEvents(
		ctx,
		topics,
		msg.ChannelID,
		channel.BotID,
	); err != nil {
		return struct{}{}
	}

	return struct{}{}
}

func (c *BusListener) subscribe(
	ctx context.Context,
	msg eventsub.EventsubSubscribeRequest,
) struct{} {
	if err := c.eventSubClient.SubscribeToEvent(
		ctx,
		msg.ConditionType,
		msg.Topic,
		msg.Version,
		msg.ChannelID,
	); err != nil {
		c.logger.Error("error subscribing to event", err)
	}

	return struct{}{}
}

func (c *BusListener) reinitChannels(
	ctx context.Context,
	_ struct{},
) struct{} {
	if err := c.eventSubClient.InitChannels(); err != nil {
		c.logger.Error("error reinit channels", err)
	}

	return struct{}{}
}
