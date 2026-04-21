package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	"channel.followed",
	"livestream.status.updated",
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

type subscribeEvent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
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
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

type subscribeRequest struct {
	BroadcasterUserID int              `json:"broadcaster_user_id"`
	Events            []subscribeEvent `json:"events"`
	Method            string           `json:"method"`
}

type subscriptionData struct {
	ID                string `json:"id"`
	Event             string `json:"event"`
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
}

type subscribeResponseItem struct {
	Name           string `json:"name"`
	Version        int    `json:"version"`
	SubscriptionID string `json:"subscription_id"`
	Error          any    `json:"error"`
}

type subscribeResponse struct {
	Data []subscribeResponseItem `json:"data"`
}

type SubscriptionInfo struct {
	SubscriptionID    string
	Event             string
	BroadcasterUserID int
	Status            string
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

type tokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        string   `json:"scope"`
	Scopes       []string `json:"scopes"`
}

func (m *SubscriptionManager) refreshBotToken(
	ctx context.Context,
	botID uuid.UUID,
	refreshToken string,
) (string, string, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", m.config.KickClientId)
	form.Set("client_secret", m.config.KickClientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://id.kick.com/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create token refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute token refresh request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read token refresh response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf(
			"kick token refresh API returned status %d: %s",
			resp.StatusCode,
			string(respBytes),
		)
	}

	var parsed tokenResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal token refresh response: %w", err)
	}

	newRefreshToken := parsed.RefreshToken
	if newRefreshToken == "" {
		newRefreshToken = refreshToken
	}

	encryptedAccessToken, err := crypto.Encrypt(parsed.AccessToken, m.config.TokensCipherKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt refreshed access token: %w", err)
	}

	encryptedRefreshToken, err := crypto.Encrypt(newRefreshToken, m.config.TokensCipherKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt refreshed refresh token: %w", err)
	}

	responseScopes := parsed.Scopes
	if len(responseScopes) == 0 && parsed.Scope != "" {
		responseScopes = strings.Fields(parsed.Scope)
	}

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
		return "", "", fmt.Errorf("failed to persist refreshed token: %w", err)
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
	reqBody := subscribeRequest{
		BroadcasterUserID: broadcasterUserID,
		Events: []subscribeEvent{
			{
				Name:    eventType,
				Version: 1,
			},
		},
		Method: "webhook",
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
		return "", &kickAPIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBytes),
		}
	}

	var subResp subscribeResponse
	if err := json.Unmarshal(respBytes, &subResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal subscribe response: %w", err)
	}

	if len(subResp.Data) == 0 {
		return "", fmt.Errorf("kick subscribe API returned empty data array")
	}

	first := subResp.Data[0]
	if first.Error != nil {
		return "", fmt.Errorf("kick subscribe API returned error for event %s: %v", first.Name, first.Error)
	}

	if first.SubscriptionID == "" {
		return "", fmt.Errorf("kick subscribe API returned empty subscription_id")
	}

	return first.SubscriptionID, nil
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

	// Build map of existing active subscriptions by event type
	existingByEvent := make(map[string][]string)
	for _, sub := range existingSubs {
		if sub.Status == "active" {
			existingByEvent[sub.Event] = append(existingByEvent[sub.Event], sub.SubscriptionID)
		}
	}

	for _, eventType := range EventTypes {
		// Check if we already have an active subscription for this event type
		if existingSubIDs, ok := existingByEvent[eventType]; ok && len(existingSubIDs) > 0 {
			// Use the first existing subscription and store it
			subID := existingSubIDs[0]
			key := redisKey(kickChannelID, eventType)
			if err := m.redis.Set(ctx, key, subID, redisTTL).Err(); err != nil {
				return fmt.Errorf("failed to store existing subscription ID for %q in Redis: %w", eventType, err)
			}

			m.logger.InfoContext(
				ctx,
				"Kick EventSub subscription already exists, reusing",
				slog.String("kick_channel_id", kickChannelID),
				slog.String("event_type", eventType),
				slog.String("subscription_id", subID),
			)

			// If there are duplicate subscriptions, clean them up
			for i := 1; i < len(existingSubIDs); i++ {
				if err := m.unsubscribe(ctx, existingSubIDs[i], broadcasterToken); err != nil {
					m.logger.WarnContext(ctx, "failed to clean up duplicate kick subscription",
						slog.String("kick_channel_id", kickChannelID),
						slog.String("event_type", eventType),
						slog.String("subscription_id", existingSubIDs[i]),
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
	u, err := url.Parse(m.apiBaseURL + "/events/subscriptions")
	if err != nil {
		return fmt.Errorf("failed to parse unsubscribe URL: %w", err)
	}

	q := u.Query()
	q.Set("id", subscriptionID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		u.String(),
		nil,
	)
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
		return &kickAPIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBytes),
		}
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

			if err := m.unsubscribe(ctx, subID, newAccessToken); err != nil {
				m.logger.WarnContext(
					ctx,
					"Failed to unsubscribe Kick EventSub after token refresh (continuing cleanup)",
					slog.String("kick_channel_id", kickChannelID),
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

type listResponse struct {
	Data []subscriptionData `json:"data"`
}

func (m *SubscriptionManager) ListSubscriptions(
	ctx context.Context,
	broadcasterToken string,
	botID uuid.UUID,
	encryptedRefreshToken string,
	broadcasterUserID int,
) ([]SubscriptionInfo, error) {
	u, err := url.Parse(m.apiBaseURL + "/events/subscriptions")
	if err != nil {
		return nil, fmt.Errorf("failed to parse list subscriptions URL: %w", err)
	}

	q := u.Query()
	q.Set("broadcaster_user_id", strconv.Itoa(broadcasterUserID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
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

	if (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) && encryptedRefreshToken != "" {
		decryptedRefreshToken, decryptErr := crypto.Decrypt(encryptedRefreshToken, m.config.TokensCipherKey)
		if decryptErr == nil {
			newAccessToken, _, refreshErr := m.refreshBotToken(ctx, botID, decryptedRefreshToken)
			if refreshErr == nil {
				return m.ListSubscriptions(ctx, newAccessToken, botID, "", broadcasterUserID)
			}
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &kickAPIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBytes),
		}
	}

	var listResp listResponse
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list subscriptions response: %w", err)
	}

	result := make([]SubscriptionInfo, len(listResp.Data))
	for i, d := range listResp.Data {
		result[i] = SubscriptionInfo{
			SubscriptionID:    d.ID,
			Event:             d.Event,
			BroadcasterUserID: d.BroadcasterUserID,
			Status:            d.Status,
			CreatedAt:         d.CreatedAt,
		}
	}

	return result, nil
}


