package manager

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/kvizyx/twitchy/eventsub"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/eventsub/internal/handler"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Manager struct {
	config             cfg.Config
	logger             logger.Logger
	gorm               *gorm.DB
	twirBus            *buscore.Bus
	conduitsRepository twitchconduits.Repository
	redSync            *redsync.Redsync
	eventsub           eventsub.EventSub
	handler            *handler.Handler

	wsCurrentSessionId *string
	currentConduit     *conduitsResponseConduit
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config             cfg.Config
	Logger             logger.Logger
	Gorm               *gorm.DB
	TwirBus            *buscore.Bus
	ConduitsRepository twitchconduits.Repository
	Redis              *goredislib.Client
	Handler            *handler.Handler
}

func NewManager(opts Opts) (*Manager, error) {
	manager := &Manager{
		config:             opts.Config,
		logger:             opts.Logger,
		gorm:               opts.Gorm,
		twirBus:            opts.TwirBus,
		conduitsRepository: opts.ConduitsRepository,
		redSync:            redsync.New(goredis.NewPool(opts.Redis)),
		eventsub:           eventsub.New(),
		handler:            opts.Handler,
		wsCurrentSessionId: nil,
		currentConduit:     nil,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := manager.createConduit(); err != nil {
					return err
				}
				go manager.startWebSocket()

				return nil
			},
		},
	)

	return manager, nil
}

func (c *Manager) SubscribeToNeededEvents(
	ctx context.Context,
	topics []model.EventsubTopic,
	broadcasterId,
	botId string,
) error {
	if c.currentConduit == nil {
		return errors.New("current conduit is not set")
	}

	var wg sync.WaitGroup
	newSubsCount := atomic.NewInt64(0)

	for _, topic := range topics {
		wg.Add(1)

		topic := topic
		go func() {
			defer wg.Done()

			err := c.SubscribeWithLimits(
				ctx,
				eventsub.EventType(topic.Topic),
				eventsub.ConduitTransport{
					Method:    "conduit",
					ConduitId: c.currentConduit.Id,
				},
				topic.Version,
				broadcasterId,
				botId,
			)

			if err != nil {
				c.logger.Error(
					"failed to subscribe to event",
					slog.Any("err", err),
					slog.Any("topic", topic.Topic),
					slog.String("version", topic.Version),
				)
			}

			newSubsCount.Inc()
		}()
	}

	wg.Wait()

	if newSubsCount.Load() > 0 {
		c.logger.Info(
			"New subscriptions created for channel",
			slog.String("channel_id", broadcasterId),
			slog.String("bot_id", botId),
			slog.Int64("count", newSubsCount.Load()),
		)
	}

	return nil
}

func (c *Manager) SubscribeToEvent(
	ctx context.Context,
	topic,
	version,
	channelId string,
) error {
	channel := model.Channels{}
	err := c.gorm.
		WithContext(ctx).
		Where(
			`"id" = ?`,
			channelId,
		).First(&channel).Error
	if err != nil {
		return err
	}

	err = c.SubscribeWithLimits(
		ctx,
		eventsub.EventType(topic),
		eventsub.ConduitTransport{
			Method:    "conduit",
			ConduitId: c.currentConduit.Id,
		},
		version,
		channel.ID,
		channel.BotID,
	)

	return err
}
