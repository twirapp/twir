package manager

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Manager struct {
	*eventsub_framework.SubClient

	config     cfg.Config
	logger     logger.Logger
	tokensGrpc tokens.TokensClient
	gorm       *gorm.DB
	tunnel     *tunnel.AppTunnel
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config     cfg.Config
	Logger     logger.Logger
	Creds      *Creds
	TokensGrpc tokens.TokensClient
	Gorm       *gorm.DB
	Tunnel     *tunnel.AppTunnel
}

func NewManager(opts Opts) (*Manager, error) {
	client := eventsub_framework.NewSubClient(opts.Creds)

	manager := &Manager{
		SubClient:  client,
		config:     opts.Config,
		logger:     opts.Logger,
		tokensGrpc: opts.TokensGrpc,
		gorm:       opts.Gorm,
		tunnel:     opts.Tunnel,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if opts.Config.AppEnv != "production" {
						if err := manager.InitChannels(); err != nil {
							panic(err)
						}
					}

					manager.SubscribeWithLimits(
						context.Background(),
						&eventsub_framework.SubRequest{
							Type: "user.authorization.revoke",
							Condition: map[string]string{
								"client_id": manager.config.TwitchClientId,
							},
							Callback: manager.tunnel.GetAddr(),
							Secret:   manager.config.TwitchClientSecret,
							Version:  "1",
						},
					)
				}()

				return nil
			},
		},
	)

	return manager, nil
}

func (c *Manager) InitChannels() error {
	if err := c.
		gorm.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&model.EventsubSubscription{}).
		Error; err != nil {
		return err
	}

	return c.populateChannels()
}

func (c *Manager) SubscribeToNeededEvents(
	ctx context.Context,
	topics []model.EventsubTopic,
	broadcasterId,
	botId string,
) error {
	var wg sync.WaitGroup
	newSubsCount := atomic.NewInt64(0)

	if err := c.unsubscribeChannel(ctx, broadcasterId); err != nil {
		c.logger.Error(
			"failed to unsubscribe from topics",
			slog.Any("err", err),
			slog.String("channel_id", broadcasterId),
		)
	}

	for _, topic := range topics {
		wg.Add(1)

		topic := topic
		go func() {
			defer wg.Done()
			condition := GetTypeCondition(topic.ConditionType, topic.Topic, broadcasterId, botId)
			if condition == nil {
				c.logger.Error(
					"failed to get condition",
					slog.String("topic", topic.Topic),
					slog.String("channel_id", broadcasterId),
					slog.String("condition_type", string(topic.ConditionType)),
				)
				return
			}

			_, err := c.SubscribeWithLimits(
				ctx,
				&eventsub_framework.SubRequest{
					Type:      topic.Topic,
					Condition: condition,
					Callback:  c.tunnel.GetAddr(),
					Secret:    c.config.TwitchClientSecret,
					Version:   topic.Version,
				},
			)

			var casterErr *eventsub_framework.TwitchError
			if err != nil && !errors.As(err, &casterErr) {
				c.logger.Error(
					"failed to subscribe to event",
					slog.Any("err", err),
					slog.Any("topic", topic.Topic),
					slog.Any("condition", condition),
					slog.String("version", topic.Version),
					slog.String("callback", c.tunnel.GetAddr()),
				)
				return
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
	conditionType,
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

	convertedCondition := model.FindEventsubCondition(conditionType)
	if conditionType == "" {
		return errors.New("condition type not found")
	}

	condition := GetTypeCondition(convertedCondition, topic, channel.ID, channel.BotID)

	if condition == nil {
		return errors.New("condition not found")
	}

	_, err = c.SubscribeWithLimits(
		ctx,
		&eventsub_framework.SubRequest{
			Type:      topic,
			Condition: condition,
			Callback:  c.tunnel.GetAddr(),
			Secret:    c.config.TwitchClientSecret,
			Version:   version,
		},
	)

	return err
}
