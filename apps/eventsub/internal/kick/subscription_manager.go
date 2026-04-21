package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/logger"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

const (
	kickAPIBase    = "https://api.kick.com/public/v1"
	redisTTL       = 25 * time.Hour
	redisKeyPrefix = "kick:sub:"
)

var EventTypes = []string{
	"chat.message.sent",
	"channel.follow",
	"stream.online",
	"stream.offline",
}

type SubscriptionManager struct {
	config          cfg.Config
	redis           *goredis.Client
	httpClient      *http.Client
	logger          *slog.Logger
	apiBaseURL      string
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
		httpClient:   http.DefaultClient,
		logger:       opts.Logger,
		apiBaseURL:   kickAPIBase,
		kickBotsRepo: opts.KickBotsRepo,
		usersRepo:    opts.UsersRepo,
	}
}

type subscribeRequest struct {
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	Type              string `json:"type"`
	Method            string `json:"method"`
	CallbackURL       string `json:"callback_url"`
}

type subscriptionData struct {
	ID                string `json:"id"`
	AppID             string `json:"app_id"`
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	Type              string `json:"type"`
	Method            string `json:"method"`
	CallbackURL       string `json:"callback_url"`
	CreatedAt         string `json:"created_at"`
}

type subscribeResponse struct {
	Data subscriptionData `json:"data"`
}

type SubscriptionInfo struct {
	ID                string
	AppID             string
	BroadcasterUserID int
	Type              string
	Method            string
	CallbackURL       string
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

func (m *SubscriptionManager) subscribe(
	ctx context.Context,
	broadcasterUserID int,
	eventType string,
	broadcasterToken string,
) (string, error) {
	reqBody := subscribeRequest{
		BroadcasterUserID: broadcasterUserID,
		Type:              eventType,
		Method:            "webhook",
		CallbackURL:       m.getWebhookCallbackURL(),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal subscribe request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		m.apiBaseURL+"/events/subscriptions",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create subscribe request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+broadcasterToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute subscribe request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read subscribe response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf(
			"kick subscribe API returned status %d: %s",
			resp.StatusCode,
			string(respBytes),
		)
	}

	var subResp subscribeResponse
	if err := json.Unmarshal(respBytes, &subResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal subscribe response: %w", err)
	}

	return subResp.Data.ID, nil
}

func (m *SubscriptionManager) SubscribeAll(
	ctx context.Context,
	kickChannelID string,
	broadcasterToken string,
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
		subID, err := m.subscribe(ctx, broadcasterUserID, eventType, broadcasterToken)
		if err != nil {
			return fmt.Errorf("failed to subscribe to %q: %w", eventType, err)
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
	reqURL := fmt.Sprintf("%s/events/subscriptions?id=%s", m.apiBaseURL, url.QueryEscape(subscriptionID))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create unsubscribe request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+broadcasterToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute unsubscribe request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"kick unsubscribe API returned status %d: %s",
			resp.StatusCode,
			string(respBytes),
		)
	}

	return nil
}

func (m *SubscriptionManager) UnsubscribeAll(ctx context.Context, kickChannelID string) error {
	kickUserUUID, err := uuid.Parse(kickChannelID)
	if err != nil {
		return fmt.Errorf("kick: parse kick channel ID %s as UUID: %w", kickChannelID, err)
	}

	kickBot, err := m.kickBotsRepo.GetByKickUserID(ctx, kickUserUUID)
	if err != nil {
		return fmt.Errorf("kick: get kick bot for channel %s: %w", kickChannelID, err)
	}

	broadcasterToken, err := crypto.Decrypt(kickBot.AccessToken, m.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("kick: decrypt access token for channel %s: %w", kickChannelID, err)
	}

	for _, eventType := range EventTypes {
		key := redisKey(kickChannelID, eventType)

		subID, err := m.redis.Get(ctx, key).Result()
		if err != nil {
			if err == goredis.Nil {
				continue
			}
			return fmt.Errorf("failed to fetch subscription ID for %q from Redis: %w", eventType, err)
		}

		if err := m.unsubscribe(ctx, subID, broadcasterToken); err != nil {
			m.logger.WarnContext(
				ctx,
				"Failed to unsubscribe Kick EventSub (continuing cleanup)",
				slog.String("kick_channel_id", kickChannelID),
				slog.String("event_type", eventType),
				logger.Error(err),
			)
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

type listResponse struct {
	Data []subscriptionData `json:"data"`
}

func (m *SubscriptionManager) ListSubscriptions(
	ctx context.Context,
	broadcasterToken string,
) ([]SubscriptionInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		m.apiBaseURL+"/events/subscriptions",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create list subscriptions request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+broadcasterToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list subscriptions request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read list subscriptions response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"kick list subscriptions API returned status %d: %s",
			resp.StatusCode,
			string(respBytes),
		)
	}

	var listResp listResponse
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list subscriptions response: %w", err)
	}

	result := make([]SubscriptionInfo, len(listResp.Data))
	for i, d := range listResp.Data {
		result[i] = SubscriptionInfo{
			ID:                d.ID,
			AppID:             d.AppID,
			BroadcasterUserID: d.BroadcasterUserID,
			Type:              d.Type,
			Method:            d.Method,
			CallbackURL:       d.CallbackURL,
			CreatedAt:         d.CreatedAt,
		}
	}

	return result, nil
}
