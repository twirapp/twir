package manager

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/google/uuid"
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
				if opts.Config.AppEnv != "production" {
					if err := manager.
						gorm.
						Session(&gorm.Session{AllowGlobalUpdate: true}).
						Delete(&model.EventsubSubscription{}).
						Error; err != nil {
						return err
					}
				}

				go func() {
					if opts.Config.AppEnv != "production" {
						twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
						if err != nil {
							panic(err)
						}

						var subscriptions []helix.EventSubSubscription
						cursor := ""
						for {
							subs, err := twitchClient.GetEventSubSubscriptions(
								&helix.EventSubSubscriptionsParams{
									After: cursor,
								},
							)
							if err != nil {
								panic(err)
							}

							subscriptions = append(subscriptions, subs.Data.EventSubSubscriptions...)

							if subs.Data.Pagination.Cursor == "" {
								break
							}

							cursor = subs.Data.Pagination.Cursor
						}

						var unsubWg sync.WaitGroup

						for _, sub := range subscriptions {
							sub := sub
							unsubWg.Add(1)
							go func() {
								defer unsubWg.Done()
								manager.Unsubscribe(ctx, sub.ID)
							}()
						}

						unsubWg.Wait()
					}

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

					manager.SubscribeWithLimits(
						requestContext,
						&eventsub_framework.SubRequest{
							Type: "user.authorization.revoke",
							Condition: map[string]string{
								"client_id": opts.Config.TwitchClientId,
							},
							Callback: opts.Tunnel.GetAddr(),
							Secret:   opts.Config.TwitchClientSecret,
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

func getTypeCondition(
	t model.EventsubConditionType,
	topic,
	channelID,
	botId string,
) map[string]string {
	switch t {
	case model.EventsubConditionTypeBroadcasterUserID:
		return map[string]string{
			"broadcaster_user_id": channelID,
		}
	case model.EventsubConditionTypeUserID:
		return map[string]string{
			"user_id": channelID,
		}
	case model.EventsubConditionTypeBroadcasterWithUserID:
		data := map[string]string{
			"broadcaster_user_id": channelID,
			"user_id":             botId,
		}
		if topic == "channel.follow" {
			data["user_id"] = channelID
		}
		return data
	case model.EventsubConditionTypeBroadcasterWithModeratorID:
		return map[string]string{
			"broadcaster_user_id": channelID,
			"moderator_user_id":   botId,
		}
	case model.EventsubConditionTypeToBroadcasterID:
		return map[string]string{
			"to_broadcaster_user_id": channelID,
		}
	default:
		return nil
	}
}

var statusesForSkip = []string{
	"enabled",
	"webhook_callback_verification_pending",
	"authorization_revoked",
	"user_removed",
	"version_removed",
}

func (c *Manager) SubscribeToNeededEvents(ctx context.Context, broadcasterId, botId string) error {
	var topics []model.EventsubTopic
	if err := c.gorm.WithContext(ctx).Find(&topics).Error; err != nil {
		return err
	}

	var existedSubscriptions []model.EventsubSubscription
	if err := c.gorm.
		WithContext(ctx).
		Where(&model.EventsubSubscription{UserID: broadcasterId}).
		Find(&existedSubscriptions).
		Error; err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, topic := range topics {
		existedSubscription, ok := lo.Find(
			existedSubscriptions,
			func(sub model.EventsubSubscription) bool {
				return sub.TopicID == topic.ID
			},
		)
		if ok && slices.Contains(statusesForSkip, existedSubscription.Status) {
			continue
		}

		wg.Add(1)

		topic := topic
		go func() {
			defer wg.Done()
			condition := getTypeCondition(topic.ConditionType, topic.Topic, broadcasterId, botId)
			if condition == nil {
				return
			}

			status, err := c.SubscribeWithLimits(
				ctx,
				&eventsub_framework.SubRequest{
					Type:      topic.Topic,
					Condition: condition,
					Callback:  c.tunnel.GetAddr(),
					Secret:    c.config.TwitchClientSecret,
					Version:   topic.Version,
				},
			)
			if err != nil {
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

			if len(status.Data) == 0 {
				c.logger.Error("failed to subscribe to event", slog.Any("status", status))
				return
			}

			sub := status.Data[0]

			if err := c.gorm.Create(
				&model.EventsubSubscription{
					ID:          uuid.MustParse(sub.ID),
					TopicID:     topic.ID,
					UserID:      broadcasterId,
					Status:      sub.Status,
					Version:     sub.Version,
					CallbackUrl: c.tunnel.GetAddr(),
				},
			).Error; err != nil {
				c.logger.Error("failed to create subscription", slog.Any("err", err))
			}
		}()
	}

	wg.Wait()

	return nil
}
