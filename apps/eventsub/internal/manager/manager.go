package manager

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
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

	wg := &sync.WaitGroup{}
	wg.Add(len(neededSubs))

	var ops uint64

	for key, value := range neededSubs {
		go func(key string, value SubRequest) {
			defer wg.Done()

			existedSub, ok := lo.Find(
				subscriptions, func(item helix.EventSubSubscription) bool {
					return item.Type == key &&
						(item.Condition.BroadcasterUserID == value.Condition["broadcaster_user_id"] ||
							item.Condition.UserID == value.Condition["user_id"])
				},
			)

			if ok && existedSub.Status == "enabled" && existedSub.Transport.Callback == c.tunnel.GetAddr() {
				return
			}

			if existedSub.Status == "authorization_revoked" {
				return
			}

			if ok {
				err = c.Unsubscribe(ctx, existedSub.ID)
				if err != nil {
					c.logger.Error(
						"Failed to unsubcribe",
						slog.String("user_id", userId),
						slog.String("key", key),
						slog.Any("err", err),
					)
					return
				}
			}

			request := eventsub_framework.SubRequest{
				Type:      key,
				Condition: value.Condition,
				Callback:  c.tunnel.GetAddr(),
				Secret:    c.config.TwitchClientSecret,
				Version:   value.Version,
			}

			retry.Do(
				func() error {
					if _, subscribeErr := c.Subscribe(ctx, &request); subscribeErr != nil {
						var e *eventsub_framework.TwitchError
						if errors.As(subscribeErr, &e) && e.Status != 409 {
							if e.Status == 429 {
								return errors.New("rate limit")
							}

							c.logger.Error(
								"Failed to subcribe",
								slog.String("user_id", userId),
								slog.String("key", key),
								slog.Any("err", e),
								slog.Int("status", e.Status),
								slog.String("message", e.Message),
								slog.String("callback", c.tunnel.GetAddr()),
							)
						} else {
							c.logger.Error(
								subscribeErr.Error(),
								slog.String("user_id", userId),
								slog.String("key", key),
							)
						}

						return nil
					}

					return nil
				},
				retry.Attempts(0),
				retry.Delay(1*time.Second),
			)

			atomic.AddUint64(&ops, 1)
		}(key, value)
	}

	wg.Wait()

	c.logger.Info(
		"Subscribed to needed events",
		slog.String("user_id", userId),
		slog.Uint64("subscriptions", ops),
	)

	return nil
}
