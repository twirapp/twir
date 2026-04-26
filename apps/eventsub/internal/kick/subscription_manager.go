package kick

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scorfly/gokick"
	goredis "github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/logger"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
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
	"livestream.status.updated",
}

var eventTypeToSubscriptionName = map[string]gokick.SubscriptionName{
	"chat.message.sent":         gokick.SubscriptionNameChatMessage,
	"channel.followed":          gokick.SubscriptionNameChannelFollow,
	"livestream.status.updated": gokick.SubscriptionNameLivestreamStatusUpdated,
}

type SubscriptionManager struct {
	config          cfg.Config
	redis           *goredis.Client
	logger          *slog.Logger
	kickBotsRepo    kickbotsrepository.Repository
	usersRepo       usersrepository.Repository
	callbackBaseURL string
}

type Opts struct {
	fx.In

	Config       cfg.Config
	Redis        *goredis.Client
	Logger       *slog.Logger
	KickBotsRepo kickbotsrepository.Repository
	UsersRepo    usersrepository.Repository
}

func New(opts Opts) *SubscriptionManager {
	return &SubscriptionManager{
		config:       opts.Config,
		redis:        opts.Redis,
		logger:       opts.Logger,
		kickBotsRepo: opts.KickBotsRepo,
		usersRepo:    opts.UsersRepo,
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

func (m *SubscriptionManager) refreshBotToken(
	ctx context.Context,
	botID uuid.UUID,
	refreshToken string,
) (string, string, error) {
	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			ClientID:     m.config.KickClientId,
			ClientSecret: m.config.KickClientSecret,
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("create kick client: %w", err)
	}

	parsed, err := client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("kick token refresh: %w", err)
	}

	newRefreshToken := parsed.RefreshToken
	if newRefreshToken == "" {
		newRefreshToken = refreshToken
	}

	encryptedAccessToken, err := crypto.Encrypt(parsed.AccessToken, m.config.TokensCipherKey)
	if err != nil {
		return "", "", fmt.Errorf("encrypt refreshed access token: %w", err)
	}

	encryptedRefreshToken, err := crypto.Encrypt(newRefreshToken, m.config.TokensCipherKey)
	if err != nil {
		return "", "", fmt.Errorf("encrypt refreshed refresh token: %w", err)
	}

	responseScopes := strings.Fields(parsed.Scope)

	_, err = m.kickBotsRepo.UpdateToken(
		ctx,
		botID,
		kickbotsrepository.UpdateTokenInput{
			AccessToken:         encryptedAccessToken,
			RefreshToken:        encryptedRefreshToken,
			Scopes:              responseScopes,
			ExpiresIn:           parsed.ExpiresIn,
			ObtainmentTimestamp: time.Now(),
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("persist refreshed token: %w", err)
	}

	m.logger.InfoContext(
		ctx,
		"kick bot token refreshed",
		slog.String("kick_bot_id", botID.String()),
	)

	return parsed.AccessToken, newRefreshToken, nil
}

func (m *SubscriptionManager) subscribe(
	ctx context.Context,
	broadcasterUserID int,
	eventType string,
	broadcasterToken string,
) (string, error) {
	subName, ok := eventTypeToSubscriptionName[eventType]
	if !ok {
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}

	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			UserAccessToken: broadcasterToken,
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
		var apiErr gokick.Error
		if errors.As(err, &apiErr) {
			return "", &kickAPIError{StatusCode: apiErr.Code(), Body: apiErr.Error()}
		}
		return "", err
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

// kickAPIError wraps HTTP status codes from Kick API calls.
type kickAPIError struct {
	StatusCode int
	Body       string
}

func (e *kickAPIError) Error() string {
	return fmt.Sprintf("kick API returned status %d: %s", e.StatusCode, e.Body)
}

func (e *kickAPIError) isAuthError() bool {
	return e.StatusCode == 401 || e.StatusCode == 403
}

func (m *SubscriptionManager) SubscribeAll(
	ctx context.Context,
	kickChannelID string,
	broadcasterToken string,
	botID uuid.UUID,
	encryptedRefreshToken string,
) error {
	kickUserUUID, err := uuid.Parse(kickChannelID)
	if err != nil {
		return fmt.Errorf("failed to parse kick channel ID %q as UUID: %w", kickChannelID, err)
	}

	user, err := m.usersRepo.GetByID(ctx, kickUserUUID.String())
	if err != nil {
		return fmt.Errorf("failed to get user for kick channel ID %q: %w", kickChannelID, err)
	}

	broadcasterUserID, err := strconv.Atoi(user.PlatformID)
	if err != nil {
		return fmt.Errorf("failed to parse kick platform user ID %q as int: %w", user.PlatformID, err)
	}

	// List existing subscriptions to avoid duplicates
	existingSubs, err := m.ListSubscriptions(ctx, broadcasterToken, botID, encryptedRefreshToken, broadcasterUserID)
	if err != nil {
		m.logger.WarnContext(ctx, "failed to list existing kick subscriptions, continuing with blind subscribe",
			slog.String("kick_channel_id", kickChannelID),
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
			key := redisKey(kickChannelID, eventType)
			if err := m.redis.Set(ctx, key, firstSub.SubscriptionID, redisTTL).Err(); err != nil {
				return fmt.Errorf("failed to store existing subscription ID for %q in Redis: %w", eventType, err)
			}

			m.logger.InfoContext(
				ctx,
				"Kick EventSub subscription already exists, reusing",
				slog.String("kick_channel_id", kickChannelID),
				slog.String("event_type", eventType),
				slog.String("subscription_id", firstSub.SubscriptionID),
			)

			for i := 1; i < len(subs); i++ {
				if err := m.unsubscribe(ctx, subs[i].SubscriptionID, broadcasterToken); err != nil {
					m.logger.WarnContext(ctx, "failed to clean up duplicate kick subscription",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						slog.String("subscription_id", subs[i].SubscriptionID),
						logger.Error(err),
					)
				}
			}

			continue
		}

		subID, err := m.subscribe(ctx, broadcasterUserID, eventType, broadcasterToken)
		if err != nil {
			var apiErr *kickAPIError
			if encryptedRefreshToken == "" || !errors.As(err, &apiErr) || !apiErr.isAuthError() {
				return fmt.Errorf("failed to subscribe to %q: %w", eventType, err)
			}

			decryptedRefreshToken, decryptErr := crypto.Decrypt(encryptedRefreshToken, m.config.TokensCipherKey)
			if decryptErr != nil {
				return fmt.Errorf("failed to subscribe to %q: %w", eventType, err)
			}

			newAccessToken, _, refreshErr := m.refreshBotToken(ctx, botID, decryptedRefreshToken)
			if refreshErr != nil {
				return fmt.Errorf("failed to subscribe to %q: %w", eventType, err)
			}

			subID, err = m.subscribe(ctx, broadcasterUserID, eventType, newAccessToken)
			if err != nil {
				return fmt.Errorf("failed to subscribe to %q after token refresh: %w", eventType, err)
			}

			broadcasterToken = newAccessToken
		}

		key := redisKey(kickChannelID, eventType)
		if err := m.redis.Set(ctx, key, subID, redisTTL).Err(); err != nil {
			return fmt.Errorf("failed to store subscription ID for %q in Redis: %w", eventType, err)
		}

		m.logger.InfoContext(
			ctx,
			"Kick EventSub subscription created",
			slog.String("kick_channel_id", kickChannelID),
			slog.String("event_type", eventType),
			slog.String("subscription_id", subID),
		)
	}

	return nil
}

func (m *SubscriptionManager) unsubscribe(
	ctx context.Context,
	subscriptionID string,
	broadcasterToken string,
) error {
	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			UserAccessToken: broadcasterToken,
		},
	)
	if err != nil {
		return fmt.Errorf("create kick client: %w", err)
	}

	_, err = client.DeleteSubscriptions(ctx, gokick.NewSubscriptionToDeleteFilter().SetIDs([]string{subscriptionID}))
	if err != nil {
		var apiErr gokick.Error
		if errors.As(err, &apiErr) {
			return &kickAPIError{StatusCode: apiErr.Code(), Body: apiErr.Error()}
		}
		return err
	}

	return nil
}

func (m *SubscriptionManager) UnsubscribeAll(
	ctx context.Context,
	kickChannelID string,
	broadcasterToken string,
	botID uuid.UUID,
	encryptedRefreshToken string,
) error {
	kickUserUUID, err := uuid.Parse(kickChannelID)
	if err != nil {
		return fmt.Errorf("failed to parse kick channel ID %q as UUID: %w", kickChannelID, err)
	}

	user, err := m.usersRepo.GetByID(ctx, kickUserUUID.String())
	if err != nil {
		return fmt.Errorf("failed to get user for kick channel ID %q: %w", kickChannelID, err)
	}

	broadcasterUserID, err := strconv.Atoi(user.PlatformID)
	if err != nil {
		return fmt.Errorf("failed to parse kick platform user ID %q as int: %w", user.PlatformID, err)
	}

	for _, eventType := range EventTypes {
		key := redisKey(kickChannelID, eventType)

		subIDs := make([]string, 0, 1)

		subID, err := m.redis.Get(ctx, key).Result()
		if err != nil {
			if err == goredis.Nil {
				subs, listErr := m.ListSubscriptions(ctx, broadcasterToken, botID, encryptedRefreshToken, broadcasterUserID)
				if listErr != nil {
					m.logger.WarnContext(
						ctx,
						"Failed to list Kick subscriptions for cleanup fallback",
						slog.String("kick_channel_id", kickChannelID),
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
			if err := m.unsubscribe(ctx, id, broadcasterToken); err != nil {
				var apiErr *kickAPIError
				shouldRefresh := encryptedRefreshToken != "" && errors.As(err, &apiErr) && apiErr.isAuthError()

				if !shouldRefresh {
					m.logger.WarnContext(
						ctx,
						"Failed to unsubscribe Kick EventSub (continuing cleanup)",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						logger.Error(err),
					)
					continue
				}

				decryptedRefreshToken, decryptErr := crypto.Decrypt(encryptedRefreshToken, m.config.TokensCipherKey)
				if decryptErr != nil {
					m.logger.WarnContext(
						ctx,
						"Failed to unsubscribe Kick EventSub (continuing cleanup)",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						logger.Error(err),
					)
					continue
				}

				newAccessToken, _, refreshErr := m.refreshBotToken(ctx, botID, decryptedRefreshToken)
				if refreshErr != nil {
					m.logger.WarnContext(
						ctx,
						"Failed to unsubscribe Kick EventSub (continuing cleanup)",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						logger.Error(err),
					)
					continue
				}

				if err := m.unsubscribe(ctx, id, newAccessToken); err != nil {
					m.logger.WarnContext(
						ctx,
						"Failed to unsubscribe Kick EventSub after token refresh (continuing cleanup)",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						logger.Error(err),
					)
				}
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
	broadcasterToken string,
	botID uuid.UUID,
	encryptedRefreshToken string,
	broadcasterUserID int,
) ([]SubscriptionInfo, error) {
	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			UserAccessToken: broadcasterToken,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create kick client: %w", err)
	}

	response, err := client.GetSubscriptions(ctx)
	if err != nil {
		var apiErr gokick.Error
		if errors.As(err, &apiErr) && (apiErr.Code() == 401 || apiErr.Code() == 403) && encryptedRefreshToken != "" {
			decryptedRefreshToken, decryptErr := crypto.Decrypt(encryptedRefreshToken, m.config.TokensCipherKey)
			if decryptErr == nil {
				newAccessToken, _, refreshErr := m.refreshBotToken(ctx, botID, decryptedRefreshToken)
				if refreshErr == nil {
					return m.ListSubscriptions(ctx, newAccessToken, botID, "", broadcasterUserID)
				}
			}
		}
		return nil, fmt.Errorf("list kick subscriptions: %w", err)
	}

	result := make([]SubscriptionInfo, 0, len(response.Result))
	for _, d := range response.Result {
		result = append(result, SubscriptionInfo{
			SubscriptionID:    d.ID,
			Event:             d.Event,
			BroadcasterUserID: d.BroadcasterUserID,
			CreatedAt:         d.CreatedAt,
		})
	}

	return result, nil
}
