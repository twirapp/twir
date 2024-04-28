package manager

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
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

type ManagerOpts struct {
	fx.In
	Lc fx.Lifecycle

	Config     cfg.Config
	Logger     logger.Logger
	Creds      *Creds
	TokensGrpc tokens.TokensClient
	Gorm       *gorm.DB
	Tunnel     *tunnel.AppTunnel
}

func NewManager(opts ManagerOpts) (*Manager, error) {
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
					requestContext := context.Background()
					var channels []model.Channels
					err := manager.gorm.Where(
						`"channels"."isEnabled" = ? AND "User"."is_banned" = ? AND "channels"."isTwitchBanned" = ?`,
						true,
						false,
						false,
					).Joins("User").Find(&channels).Error
					if err != nil {
						panic(err)
					}

					for _, channel := range channels {
						err = manager.SubscribeToNeededEvents(requestContext, channel.ID, channel.BotID)
						if err != nil {
							continue
						}
					}
				}()

				return nil
			},
		},
	)

	return manager, nil
}

type SubRequest struct {
	Version   string
	Condition map[string]string
}

func (c *Manager) SubscribeToNeededEvents(ctx context.Context, userId, botId string) error {
	channelCondition := map[string]string{
		"broadcaster_user_id": userId,
	}
	userCondition := map[string]string{
		"user_id": userId,
	}

	channelConditionWithBotId := map[string]string{
		"broadcaster_user_id": userId,
		"user_id":             botId,
	}

	channelConditionWithModeratorId := map[string]string{
		"broadcaster_user_id": userId,
		"moderator_user_id":   botId,
	}

	neededSubs := map[string]SubRequest{
		"channel.update": {
			Version:   "2",
			Condition: channelCondition,
		},
		"stream.online": {
			Version:   "1",
			Condition: channelCondition,
		},
		"stream.offline": {
			Version:   "1",
			Condition: channelCondition,
		},
		"user.update": {
			Condition: userCondition,
			Version:   "1",
		},
		"channel.follow": {
			Version: "2",
			Condition: map[string]string{
				"broadcaster_user_id": userId,
				"moderator_user_id":   userId,
			},
		},
		"channel.moderator.add": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.moderator.remove": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.channel_points_custom_reward_redemption.add": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.channel_points_custom_reward_redemption.update": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.begin": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.progress": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.poll.end": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.begin": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.lock": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.progress": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.prediction.end": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.ban": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.subscribe": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.subscription.gift": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.subscription.message": {
			Version:   "1",
			Condition: channelCondition,
		},
		"channel.raid": {
			Version: "1",
			Condition: map[string]string{
				"to_broadcaster_user_id": userId,
			},
		},
		"channel.chat.clear": {
			Version:   "1",
			Condition: channelConditionWithBotId,
		},
		"channel.chat.clear_user_messages": {
			Version:   "1",
			Condition: channelConditionWithBotId,
		},
		"channel.chat.message_delete": {
			Version:   "1",
			Condition: channelConditionWithBotId,
		},
		"channel.chat.notification": {
			Version:   "1",
			Condition: channelConditionWithBotId,
		},
		"channel.chat.message": {
			Version:   "1",
			Condition: channelConditionWithBotId,
		},
		"channel.unban_request.create": {
			Version:   "1",
			Condition: channelConditionWithModeratorId,
		},
		"channel.unban_request.resolve": {
			Version:   "1",
			Condition: channelConditionWithModeratorId,
		},
	}

	twitchClient, err := twitch.NewAppClient(c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	var subscriptions []helix.EventSubSubscription
	cursor := ""
	for {
		subs, err := twitchClient.GetEventSubSubscriptions(
			&helix.EventSubSubscriptionsParams{
				UserID: userId,
				After:  cursor,
			},
		)
		if err != nil {
			return err
		}

		subscriptions = append(subscriptions, subs.Data.EventSubSubscriptions...)

		if subs.Data.Pagination.Cursor == "" {
			break
		}

		cursor = subs.Data.Pagination.Cursor
	}

	if c.config.AppEnv != "production" {
		for _, sub := range subscriptions {
			c.Unsubscribe(ctx, sub.ID)
		}

		for key, value := range neededSubs {
			c.Subscribe(
				ctx, &eventsub_framework.SubRequest{
					Type:      key,
					Condition: value.Condition,
					Callback:  c.tunnel.GetAddr(),
					Secret:    c.config.TwitchClientSecret,
					Version:   value.Version,
				},
			)

			c.logger.Info(
				"Subscribed",
				slog.String("type", key),
				slog.String("user_id", userId),
			)
		}
	} else {
		for key, value := range neededSubs {
			for _, sub := range subscriptions {
				if sub.Type != key {
					continue
				}

				if sub.Status == "notification_failures_exceeded" {
					c.logger.Info(
						"Notification failures exceeded, resubscribing",
						slog.String("type", key),
						slog.String("user_id", userId),
					)

					if err := retry.Do(
						func() error {
							_, subscribeErr := c.Subscribe(
								ctx, &eventsub_framework.SubRequest{
									Type:      key,
									Condition: value.Condition,
									Callback:  c.tunnel.GetAddr(),
									Secret:    c.config.TwitchClientSecret,
									Version:   value.Version,
								},
							)

							return subscribeErr
						},
						retry.Attempts(0),
						retry.Delay(1*time.Second),
						retry.RetryIf(
							func(err error) bool {
								var e *eventsub_framework.TwitchError
								if errors.As(err, &e) && e.Status != 409 {
									if e.Status == 429 {
										return true
									}
								}

								return false
							},
						),
					); err != nil {
						c.logger.Error(
							"Failed to resubscribe",
							slog.Any("err", err),
							slog.String("type", key),
							slog.String("user_id", userId),
						)
					}
				}
			}

			_, isExists := lo.Find(
				subscriptions, func(item helix.EventSubSubscription) bool {
					return item.Type == key
				},
			)

			if !isExists {
				c.logger.Info(
					"Subscription not found, resubscribing",
					slog.String("type", key),
					slog.String("user_id", userId),
				)

				if _, err := c.Subscribe(
					ctx, &eventsub_framework.SubRequest{
						Type:      key,
						Condition: value.Condition,
						Callback:  c.tunnel.GetAddr(),
						Secret:    c.config.TwitchClientSecret,
						Version:   value.Version,
					},
				); err != nil {
					c.logger.Error(
						"Failed to resubscribe",
						slog.Any("err", err),
						slog.String("type", key),
						slog.String("user_id", userId),
					)
				}
			}
		}
	}

	return nil
}
