package kick

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scorfly/gokick"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	kick_bots_entity "github.com/twirapp/twir/libs/entities/kick_bot"
	"github.com/twirapp/twir/libs/repositories/kick_bots"
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

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: decryptedAccessToken,
		HTTPClient:      c.httpClient,
		ClientID:        c.config.KickClientId,
		ClientSecret:    c.config.KickClientSecret,
	})
	if err != nil {
		return fmt.Errorf("create gokick client: %w", err)
	}

	_, err = kickClient.SendChatMessage(ctx, &broadcasterUserID, text, nil, gokick.MessageTypeUser)
	if err != nil {
		var apiErr gokick.Error
		if errors.As(err, &apiErr) {
			switch apiErr.Code() {
			case http.StatusUnauthorized:
				if !allowRefresh {
					return fmt.Errorf("kick chat request failed with status %d: %w", apiErr.Code(), apiErr)
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
					slog.Int("status_code", apiErr.Code()),
					slog.Any("error", err),
				)
				return nil
			case http.StatusForbidden:
				c.logger.WarnContext(
					ctx,
					"kick chat forbidden",
					slog.String("broadcaster_kick_id", broadcasterKickID),
					slog.Int("broadcaster_user_id", broadcasterUserID),
					slog.String("bot_id", bot.ID),
					slog.String("bot_kick_login", bot.KickUserLogin),
					slog.String("bot_kick_user_id", bot.KickUserID.String()),
					slog.Any("bot_scopes", bot.Scopes),
					slog.Int("status_code", apiErr.Code()),
					slog.Any("error", err),
				)
				return fmt.Errorf("kick chat request failed with status %d: %w", apiErr.Code(), apiErr)
			default:
				return fmt.Errorf("kick chat request failed with status %d: %w", apiErr.Code(), apiErr)
			}
		}
		return fmt.Errorf("send kick chat message: %w", err)
	}

	return nil
}

func (c *ChatClient) refreshBotToken(ctx context.Context, bot *kick_bots_entity.KickBot) error {
	decryptedRefreshToken, err := crypto.Decrypt(bot.RefreshToken, c.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("decrypt kick bot refresh token: %w", err)
	}

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		HTTPClient:   c.httpClient,
		ClientID:     c.config.KickClientId,
		ClientSecret: c.config.KickClientSecret,
	})
	if err != nil {
		return fmt.Errorf("create gokick client: %w", err)
	}

	tokenResp, err := kickClient.RefreshToken(ctx, decryptedRefreshToken)
	if err != nil {
		return fmt.Errorf("kick token refresh failed: %w", err)
	}

	responseScopes := bot.Scopes
	if tokenResp.Scope != "" {
		responseScopes = strings.Fields(tokenResp.Scope)
	}

	refreshToken := tokenResp.RefreshToken
	if refreshToken == "" {
		refreshToken = decryptedRefreshToken
	}

	encryptedAccessToken, err := crypto.Encrypt(tokenResp.AccessToken, c.config.TokensCipherKey)
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
			ExpiresIn:           tokenResp.ExpiresIn,
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
		slog.String("kick_user_login", bot.KickUserLogin),
		slog.Any("scopes", responseScopes),
		slog.Int("expires_in", tokenResp.ExpiresIn),
	)

	return nil
}

func splitMessage(text string) []string {
	normalizedText := strings.ReplaceAll(text, "\n", " ")
	if normalizedText == "" {
		return nil
	}

	if err := gokick.ValidateChatMessageContent(normalizedText); err == nil {
		return []string{normalizedText}
	}

	runes := []rune(normalizedText)
	parts := make([]string, 0, (len(runes)/gokick.ChatMessageContentMaxRunes)+1)
	for len(runes) > 0 {
		end := gokick.ChatMessageContentMaxRunes
		if len(runes) < gokick.ChatMessageContentMaxRunes {
			end = len(runes)
		}

		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}

	return parts
}
