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

	"github.com/scorfly/gokick"
	cfg "github.com/twirapp/twir/libs/config"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	kickoauth "github.com/twirapp/twir/libs/oauth/kick"
)

type ChatClient struct {
	config          cfg.Config
	httpClient      *http.Client
	logger          *slog.Logger
	botTokens       kickoauth.DefaultBotTokenSource
}

func NewChatClient(botTokens kickoauth.DefaultBotTokenSource, config cfg.Config) *ChatClient {
	return &ChatClient{
		config:          config,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		logger:          slog.Default(),
		botTokens:       botTokens,
	}
}

func (c *ChatClient) SendMessage(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	text string,
	replyToMessageID string,
) error {
	broadcasterKickID := binding.PlatformChannelID
	broadcasterUserID, err := strconv.Atoi(broadcasterKickID)
	if err != nil {
		return fmt.Errorf("parse broadcaster kick id: %w", err)
	}

	for _, part := range splitMessage(text) {
		err = c.sendMessagePart(ctx, broadcasterUserID, broadcasterKickID, part, replyToMessageID)
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
	replyToMessageID string,
) error {
	token, err := c.botTokens.Token(ctx)
	if err != nil {
		return fmt.Errorf("request kick bot token: %w", err)
	}

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: token.AccessToken,
		HTTPClient:      c.httpClient,
		ClientID:        c.config.KickClientId,
		ClientSecret:    c.config.KickClientSecret,
	})
	if err != nil {
		return fmt.Errorf("create gokick client: %w", err)
	}

	var replyToPtr *string
	if replyToMessageID != "" {
		replyToPtr = &replyToMessageID
	}

	_, err = kickClient.SendChatMessage(ctx, &broadcasterUserID, text, replyToPtr, gokick.MessageTypeUser)
	if err != nil {
		var apiErr gokick.Error
		if errors.As(err, &apiErr) {
			switch apiErr.Code() {
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
					slog.Any("bot_scopes", token.Scopes),
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

func splitMessage(text string) []string {
	normalizedText := strings.ReplaceAll(text, "\n", " ")
	if normalizedText == "" {
		return nil
	}

	if err := gokick.ValidateChatMessageContent(normalizedText); err == nil {
		return []string{normalizedText}
	}

	parts := make([]string, 0, (len(normalizedText)/gokick.ChatMessageContentMaxRunes)+1)
	for len(normalizedText) > 0 {
		if len(normalizedText) <= gokick.ChatMessageContentMaxRunes {
			parts = append(parts, normalizedText)
			break
		}

		end := 0
		for i := range normalizedText {
			if i == 0 {
				continue
			}
			if i > gokick.ChatMessageContentMaxRunes {
				break
			}
			end = i
		}

		if end == 0 {
			end = gokick.ChatMessageContentMaxRunes
		}

		parts = append(parts, normalizedText[:end])
		normalizedText = normalizedText[end:]
	}

	return parts
}
