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
	"strings"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	kick_bots_entity "github.com/twirapp/twir/libs/entities/kick_bot"
	"github.com/twirapp/twir/libs/repositories/kick_bots"
)

const (
	kickAPIBaseURL       = "https://api.kick.com/public/v1"
	kickOAuthBaseURL     = "https://id.kick.com"
	maxKickMessageLength = 500
)

type ChatClient struct {
	repo       kick_bots.Repository
	config     cfg.Config
	httpClient *http.Client
	logger     *slog.Logger
}

func NewChatClient(repo kick_bots.Repository, config cfg.Config) *ChatClient {
	return &ChatClient{
		repo:       repo,
		config:     config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     slog.Default(),
	}
}

func (c *ChatClient) SendMessage(ctx context.Context, broadcasterKickID string, text string) error {
	broadcasterUserID, err := strconv.Atoi(broadcasterKickID)
	if err != nil {
		return fmt.Errorf("parse broadcaster kick id: %w", err)
	}

	bot, err := c.repo.GetDefault(ctx)
	if err != nil {
		return fmt.Errorf("get default kick bot: %w", err)
	}

	for _, part := range splitMessage(text) {
		err = c.sendMessagePart(ctx, broadcasterUserID, broadcasterKickID, part, &bot, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ChatClient) sendMessagePart(
	ctx context.Context,
	broadcasterUserID int,
	broadcasterKickID string,
	text string,
	bot *kick_bots_entity.KickBot,
	allowRefresh bool,
) error {
	decryptedAccessToken, err := crypto.Decrypt(bot.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("decrypt kick bot access token: %w", err)
	}

	statusCode, responseBody, err := c.doSendMessageRequest(ctx, broadcasterUserID, text, decryptedAccessToken)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusUnauthorized:
		if !allowRefresh {
			return fmt.Errorf("kick chat request failed with status %d: %s", statusCode, responseBody)
		}

		if err := c.refreshBotToken(ctx, bot); err != nil {
			return fmt.Errorf("refresh kick bot token: %w", err)
		}

		return c.sendMessagePart(ctx, broadcasterUserID, broadcasterKickID, text, bot, false)
	case http.StatusTooManyRequests:
		c.logger.WarnContext(
			ctx,
			"kick chat rate limited, dropping message",
			slog.String("broadcaster_kick_id", broadcasterKickID),
			slog.Int("status_code", statusCode),
			slog.String("response", responseBody),
		)
		return nil
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("kick chat request failed with status %d: %s", statusCode, responseBody)
	}

	return nil
}

func (c *ChatClient) doSendMessageRequest(
	ctx context.Context,
	broadcasterUserID int,
	text string,
	accessToken string,
) (int, string, error) {
	body, err := json.Marshal(kickChatRequest{
		BroadcasterUserID:      broadcasterUserID,
		Content:                text,
		Type:                   "message",
		ReplyToOriginalMessage: nil,
	})
	if err != nil {
		return 0, "", fmt.Errorf("marshal kick chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		kickAPIBaseURL+"/chat",
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, "", fmt.Errorf("create kick chat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("do kick chat request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("read kick chat response: %w", err)
	}

	return resp.StatusCode, string(responseBody), nil
}

func (c *ChatClient) refreshBotToken(ctx context.Context, bot *kick_bots_entity.KickBot) error {
	decryptedRefreshToken, err := crypto.Decrypt(bot.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("decrypt kick bot refresh token: %w", err)
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", decryptedRefreshToken)
	form.Set("client_id", c.config.KickClientId)
	form.Set("client_secret", c.config.KickClientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		kickOAuthBaseURL+"/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("create kick token refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do kick token refresh request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read kick token refresh response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("kick token refresh failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	var parsed kickTokenResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return fmt.Errorf("decode kick token refresh response: %w", err)
	}

	responseScopes := parsed.Scopes
	if len(responseScopes) == 0 && parsed.Scope != "" {
		responseScopes = strings.Fields(parsed.Scope)
	}
	if len(responseScopes) == 0 {
		responseScopes = bot.Scopes
	}

	refreshToken := parsed.RefreshToken
	if refreshToken == "" {
		refreshToken = decryptedRefreshToken
	}

	encryptedAccessToken, err := crypto.Encrypt(parsed.AccessToken, c.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt kick bot access token: %w", err)
	}

	encryptedRefreshToken, err := crypto.Encrypt(refreshToken, c.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt kick bot refresh token: %w", err)
	}

	botID, err := uuid.Parse(bot.ID)
	if err != nil {
		return fmt.Errorf("parse kick bot id: %w", err)
	}

	updatedBot, err := c.repo.UpdateToken(
		ctx,
		botID,
		kick_bots.UpdateTokenInput{
			AccessToken:         encryptedAccessToken,
			RefreshToken:        encryptedRefreshToken,
			Scopes:              responseScopes,
			ExpiresIn:           parsed.ExpiresIn,
			ObtainmentTimestamp: time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("persist kick bot token: %w", err)
	}

	*bot = updatedBot

	c.logger.InfoContext(
		ctx,
		"kick bot token refreshed",
		slog.String("kick_bot_id", bot.ID),
		slog.String("kick_user_id", bot.KickUserID.String()),
	)

	return nil
}

type kickChatRequest struct {
	BroadcasterUserID      int     `json:"broadcaster_user_id"`
	Content                string  `json:"content"`
	Type                   string  `json:"type"`
	ReplyToOriginalMessage *string `json:"reply_to_original_message_id"`
}

type kickTokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        string   `json:"scope"`
	Scopes       []string `json:"scopes"`
}

func splitMessage(text string) []string {
	normalizedText := strings.ReplaceAll(text, "\n", " ")
	if normalizedText == "" {
		return nil
	}

	runes := []rune(normalizedText)
	if len(runes) <= maxKickMessageLength {
		return []string{normalizedText}
	}

	parts := make([]string, 0, (len(runes)/maxKickMessageLength)+1)
	for len(runes) > 0 {
		end := maxKickMessageLength
		if len(runes) < maxKickMessageLength {
			end = len(runes)
		}

		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}

	return parts
}
