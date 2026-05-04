package kick

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/scorfly/gokick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	bustokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

const (
	redisTTL       = 25 * time.Hour
	redisKeyPrefix = "kick:sub:"
)

var EventTypes = []string{
	"chat.message.sent",
	"channel.followed",
	"channel.subscription.new",
	"channel.subscription.renewal",
	"channel.subscription.gifts",
	"channel.reward.redemption.updated",
	"livestream.status.updated",
	"livestream.metadata.updated",
	"moderation.banned",
}

var eventTypeToSubscriptionName = map[string]gokick.SubscriptionName{
	"chat.message.sent":                 gokick.SubscriptionNameChatMessage,
	"channel.followed":                  gokick.SubscriptionNameChannelFollow,
	"channel.subscription.new":          gokick.SubscriptionNameChannelSubscriptionCreated,
	"channel.subscription.renewal":      gokick.SubscriptionNameChannelSubscriptionRenewal,
	"channel.subscription.gifts":        gokick.SubscriptionNameChannelSubscriptionGifts,
	"channel.reward.redemption.updated": gokick.SubscriptionNameChannelRewardRedemptionUpdated,
	"livestream.status.updated":         gokick.SubscriptionNameLivestreamStatusUpdated,
	"livestream.metadata.updated":       gokick.SubscriptionNameLivestreamMetadataUpdated,
	"moderation.banned":                 gokick.SubscriptionNameModerationBanned,
}

type SubscriptionManager struct {
	config          cfg.Config
	redis           *goredis.Client
	logger          *slog.Logger
	usersRepo       usersrepository.Repository
	twirBus         *buscore.Bus
	callbackBaseURL string
}

type Opts struct {
	fx.In

	Config    cfg.Config
	Redis     *goredis.Client
	Logger    *slog.Logger
	TwirBus   *buscore.Bus
	UsersRepo usersrepository.Repository
}

func New(opts Opts) *SubscriptionManager {
	return &SubscriptionManager{
		config:    opts.Config,
		redis:     opts.Redis,
		logger:    opts.Logger,
		twirBus:   opts.TwirBus,
		usersRepo: opts.UsersRepo,
	}
}

type SubscriptionInfo struct {
	SubscriptionID    string
	Event             string
	BroadcasterUserID int
	CreatedAt         string
}

func (m *SubscriptionManager) Name() string {
	return "kick"
}

func (m *SubscriptionManager) SetCallbackBaseURL(baseURL string) {
	m.callbackBaseURL = baseURL
}

func redisKey(kickChannelID, eventType string) string {
	return redisKeyPrefix + kickChannelID + ":" + eventType
}

func (m *SubscriptionManager) getWebhookCallbackURL() string {
	baseURL := m.callbackBaseURL
	if baseURL == "" {
		baseURL = m.config.SiteBaseUrl
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL + "/webhook/kick"
	}
	return u.JoinPath("webhook", "kick").String()
}

func (m *SubscriptionManager) getAppAccessToken(ctx context.Context) (string, error) {
	resp, err := m.twirBus.Tokens.RequestAppToken.Request(
		ctx,
		bustokens.GetAppTokenRequest{Platform: platformentity.PlatformKick},
	)
	if err != nil {
		return "", fmt.Errorf("request kick app token: %w", err)
	}

	return resp.Data.AccessToken, nil
}

func (m *SubscriptionManager) subscribe(
	ctx context.Context,
	broadcasterUserID int,
	eventType string,
) (string, error) {
	subName, ok := eventTypeToSubscriptionName[eventType]
	if !ok {
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}

	appToken, err := m.getAppAccessToken(ctx)
	if err != nil {
		return "", err
	}

	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			AppAccessToken: appToken,
		},
	)
	if err != nil {
		return "", fmt.Errorf("create kick client: %w", err)
	}

	subscriptions := []gokick.SubscriptionRequest{
		{Name: subName, Version: 1},
	}

	response, err := client.CreateSubscriptions(
		ctx,
		gokick.SubscriptionMethodWebhook,
		subscriptions,
		&broadcasterUserID,
	)
	if err != nil {
		return "", fmt.Errorf("kick create subscription: %w", err)
	}

	if len(response.Result) == 0 {
		return "", fmt.Errorf("kick subscribe API returned empty result")
	}

	first := response.Result[0]
	if first.Error != "" {
		return "", fmt.Errorf("kick subscribe API returned error for event %s: %s", first.Name, first.Error)
	}

	if first.SubscriptionID == "" {
		return "", fmt.Errorf("kick subscribe API returned empty subscription_id")
	}

	return first.SubscriptionID, nil
}

func (m *SubscriptionManager) unsubscribe(
	ctx context.Context,
	subscriptionID string,
) error {
	appToken, err := m.getAppAccessToken(ctx)
	if err != nil {
		return err
	}

	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			AppAccessToken: appToken,
		},
	)
	if err != nil {
		return fmt.Errorf("create kick client: %w", err)
	}

	_, err = client.DeleteSubscriptions(ctx, gokick.NewSubscriptionToDeleteFilter().SetIDs([]string{subscriptionID}))
	if err != nil {
		return fmt.Errorf("kick delete subscription: %w", err)
	}

	return nil
}

func (m *SubscriptionManager) SubscribeAll(
	ctx context.Context,
	kickChannelID uuid.UUID,
) error {
	kickChannelIDStr := kickChannelID.String()

	user, err := m.usersRepo.GetByID(ctx, kickChannelID)
	if err != nil {
		return fmt.Errorf("failed to get user for kick channel ID %q: %w", kickChannelIDStr, err)
	}

	broadcasterUserID, err := strconv.Atoi(user.PlatformID)
	if err != nil {
		return fmt.Errorf("failed to parse kick platform user ID %q as int: %w", user.PlatformID, err)
	}

	existingSubs, err := m.ListSubscriptions(ctx, broadcasterUserID)
	if err != nil {
		m.logger.WarnContext(ctx, "failed to list existing kick subscriptions, continuing with blind subscribe",
			slog.String("kick_channel_id", kickChannelIDStr),
			logger.Error(err),
		)
		existingSubs = nil
	}

	existingByEvent := make(map[string][]SubscriptionInfo)
	for _, sub := range existingSubs {
		existingByEvent[sub.Event] = append(existingByEvent[sub.Event], sub)
	}

	for _, eventType := range EventTypes {
		subs, ok := existingByEvent[eventType]
		if ok && len(subs) > 0 {
			firstSub := &subs[0]
			key := redisKey(kickChannelIDStr, eventType)
			if err := m.redis.Set(ctx, key, firstSub.SubscriptionID, redisTTL).Err(); err != nil {
				return fmt.Errorf("failed to store existing subscription ID for %q in Redis: %w", eventType, err)
			}

			m.logger.InfoContext(
				ctx,
				"Kick EventSub subscription already exists, reusing",
				slog.String("kick_channel_id", kickChannelIDStr),
				slog.String("event_type", eventType),
				slog.String("subscription_id", firstSub.SubscriptionID),
			)

			for i := 1; i < len(subs); i++ {
				if err := m.unsubscribe(ctx, subs[i].SubscriptionID); err != nil {
					m.logger.WarnContext(ctx, "failed to clean up duplicate kick subscription",
						slog.String("kick_channel_id", kickChannelIDStr),
						slog.String("event_type", eventType),
						slog.String("subscription_id", subs[i].SubscriptionID),
						logger.Error(err),
					)
				}
			}

			continue
		}

		subID, err := m.subscribe(ctx, broadcasterUserID, eventType)
		if err != nil {
			return fmt.Errorf("failed to subscribe to %q: %w", eventType, err)
		}

		key := redisKey(kickChannelIDStr, eventType)
		if err := m.redis.Set(ctx, key, subID, redisTTL).Err(); err != nil {
			return fmt.Errorf("failed to store subscription ID for %q in Redis: %w", eventType, err)
		}

		m.logger.InfoContext(
			ctx,
			"Kick EventSub subscription created",
			slog.String("kick_channel_id", kickChannelIDStr),
			slog.String("event_type", eventType),
			slog.String("subscription_id", subID),
		)
	}

	return nil
}

func (m *SubscriptionManager) UnsubscribeAll(
	ctx context.Context,
	kickChannelID uuid.UUID,
) error {
	kickChannelIDStr := kickChannelID.String()

	user, err := m.usersRepo.GetByID(ctx, kickChannelID)
	if err != nil {
		return fmt.Errorf("failed to get user for kick channel ID %q: %w", kickChannelIDStr, err)
	}

	broadcasterUserID, err := strconv.Atoi(user.PlatformID)
	if err != nil {
		return fmt.Errorf("failed to parse kick platform user ID %q as int: %w", user.PlatformID, err)
	}

	for _, eventType := range EventTypes {
		key := redisKey(kickChannelIDStr, eventType)

		subIDs := make([]string, 0, 1)

		subID, err := m.redis.Get(ctx, key).Result()
		if err != nil {
			if err == goredis.Nil {
				subs, listErr := m.ListSubscriptions(ctx, broadcasterUserID)
				if listErr != nil {
				m.logger.WarnContext(
					ctx,
					"Failed to list Kick subscriptions for cleanup fallback",
					slog.String("kick_channel_id", kickChannelIDStr),
					slog.String("event_type", eventType),
					logger.Error(listErr),
				)
					continue
				}

				for _, sub := range subs {
					if sub.Event == eventType {
						subIDs = append(subIDs, sub.SubscriptionID)
					}
				}
				if len(subIDs) == 0 {
					continue
				}
			} else {
				return fmt.Errorf("failed to fetch subscription ID for %q from Redis: %w", eventType, err)
			}
		} else {
			subIDs = append(subIDs, subID)
		}

		for _, id := range subIDs {
			if err := m.unsubscribe(ctx, id); err != nil {
				m.logger.WarnContext(
					ctx,
					"Failed to unsubscribe Kick EventSub (continuing cleanup)",
					slog.String("kick_channel_id", kickChannelIDStr),
					slog.String("event_type", eventType),
					logger.Error(err),
				)
			}
		}

		if err := m.redis.Del(ctx, key).Err(); err != nil {
			m.logger.WarnContext(
				ctx,
				"Failed to delete Kick subscription ID from Redis",
				slog.String("key", key),
				logger.Error(err),
			)
		}
	}

	return nil
}

func (m *SubscriptionManager) ListSubscriptions(
	ctx context.Context,
	broadcasterUserID int,
) ([]SubscriptionInfo, error) {
	appToken, err := m.getAppAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			AppAccessToken: appToken,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create kick client: %w", err)
	}

	response, err := client.GetSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list kick subscriptions: %w", err)
	}

	result := make([]SubscriptionInfo, 0, len(response.Result))
	for _, d := range response.Result {
		if d.BroadcasterUserID == broadcasterUserID {
			result = append(result, SubscriptionInfo{
				SubscriptionID:    d.ID,
				Event:             d.Event,
				BroadcasterUserID: d.BroadcasterUserID,
				CreatedAt:         d.CreatedAt,
			})
		}
	}

	return result, nil
}
