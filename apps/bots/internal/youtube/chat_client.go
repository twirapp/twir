package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

const (
	youtubeAPIBaseURL         = "https://www.googleapis.com/youtube/v3"
	liveChatIDCachePrefix     = "youtube:livechat_id:"
	liveChatIDCacheTTL        = 5 * time.Minute
	noActiveBroadcastCacheTTL = time.Minute
	chatMessageMaxRunes       = 200
)

var errMessageDropped = errors.New("youtube chat message dropped")

type botTokenRequester interface {
	Request(context.Context, buscoretokens.GetBotTokenRequest) (*buscore.QueueResponse[buscoretokens.TokenResponse], error)
}

type userTokenRequester interface {
	Request(context.Context, buscoretokens.GetUserTokenRequest) (*buscore.QueueResponse[buscoretokens.TokenResponse], error)
}

type liveChatIDCache interface {
	Get(context.Context, string) kv.Valuer
	Set(context.Context, string, any, ...kvoptions.Option) error
	Delete(context.Context, string) error
}

type ChatClient struct {
	apiBaseURL       string
	httpClient       *http.Client
	logger           *slog.Logger
	liveChatIDCache  liveChatIDCache
	requestBotToken  botTokenRequester
	requestUserToken userTokenRequester
}

func NewChatClient(twirBus *buscore.Bus, _ cfg.Config, liveChatIDCache kv.KV) *ChatClient {
	return &ChatClient{
		apiBaseURL:       youtubeAPIBaseURL,
		httpClient:       &http.Client{Timeout: 10 * time.Second},
		logger:           slog.Default(),
		liveChatIDCache:  liveChatIDCache,
		requestBotToken:  twirBus.Tokens.RequestBotToken,
		requestUserToken: twirBus.Tokens.RequestUserToken,
	}
}

func (c *ChatClient) SendMessage(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	text string,
) error {
	liveChatID, err := c.resolveLiveChatID(ctx, binding)
	if errors.Is(err, errMessageDropped) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, part := range splitMessage(text) {
		if err := c.sendMessagePart(ctx, binding, liveChatID, part); err != nil {
			return err
		}
	}

	return nil
}

func (c *ChatClient) resolveLiveChatID(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
) (string, error) {
	cacheKey := liveChatIDCachePrefix + binding.ChannelID.String()
	cachedLiveChatID, err := c.liveChatIDCache.Get(ctx, cacheKey).String()
	if err == nil {
		if cachedLiveChatID == "" {
			return "", noActiveBroadcastError(binding)
		}
		return cachedLiveChatID, nil
	}
	if !errors.Is(err, kv.ErrKeyNil) {
		c.logger.WarnContext(
			ctx,
			"youtube live chat cache lookup failed",
			slog.String("channel_id", binding.ChannelID.String()),
			slog.Any("error", err),
		)
	}

	ownerTokenResponse, err := c.requestUserToken.Request(
		ctx,
		buscoretokens.GetUserTokenRequest{UserId: binding.UserID},
	)
	if err != nil {
		return "", fmt.Errorf("request YouTube binding owner token: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.endpoint("/liveBroadcasts?part=snippet&mine=true&broadcastStatus=active"),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("create YouTube live broadcasts request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+ownerTokenResponse.Data.AccessToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("request active YouTube live broadcast: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusTooManyRequests:
		c.logger.WarnContext(
			ctx,
			"youtube live broadcasts rate limited, dropping message",
			slog.String("channel_id", binding.ChannelID.String()),
			slog.Int("status_code", response.StatusCode),
		)
		return "", errMessageDropped
	case http.StatusForbidden:
		c.logger.WarnContext(
			ctx,
			"youtube live broadcasts forbidden",
			slog.String("channel_id", binding.ChannelID.String()),
			slog.Int("status_code", response.StatusCode),
		)
		return "", fmt.Errorf("request active YouTube live broadcast: %w", youtubeAPIStatusError{statusCode: response.StatusCode})
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("request active YouTube live broadcast: %w", youtubeAPIStatusError{statusCode: response.StatusCode})
	}

	var broadcastsResponse liveBroadcastsResponse
	if err := json.NewDecoder(response.Body).Decode(&broadcastsResponse); err != nil {
		return "", fmt.Errorf("decode active YouTube live broadcasts response: %w", err)
	}

	if len(broadcastsResponse.Items) == 0 || broadcastsResponse.Items[0].Snippet.LiveChatID == "" {
		c.cacheLiveChatID(ctx, cacheKey, "", noActiveBroadcastCacheTTL)
		return "", noActiveBroadcastError(binding)
	}

	liveChatID := broadcastsResponse.Items[0].Snippet.LiveChatID
	c.cacheLiveChatID(ctx, cacheKey, liveChatID, liveChatIDCacheTTL)
	return liveChatID, nil
}

func (c *ChatClient) sendMessagePart(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	liveChatID string,
	text string,
) error {
	botTokenResponse, err := c.requestBotToken.Request(
		ctx,
		buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformYouTube},
	)
	if err != nil {
		return fmt.Errorf("request YouTube bot token: %w", err)
	}

	body, err := json.Marshal(liveChatMessageInsertRequest{
		Snippet: liveChatMessageSnippet{
			LiveChatID: liveChatID,
			Type:       "textMessageEvent",
			TextMessageDetails: liveChatTextMessageDetails{
				MessageText: text,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("marshal YouTube live chat message: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint("/liveChat/messages?part=snippet"),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create YouTube live chat message request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+botTokenResponse.Data.AccessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send YouTube live chat message: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusTooManyRequests:
		c.logger.WarnContext(
			ctx,
			"youtube live chat rate limited, dropping message",
			slog.String("channel_id", binding.ChannelID.String()),
			slog.Int("status_code", response.StatusCode),
		)
		return nil
	case http.StatusForbidden, http.StatusNotFound:
		c.invalidateLiveChatID(ctx, binding.ChannelID.String())
		c.logger.WarnContext(
			ctx,
			"youtube live chat forbidden; chat may be disabled or bot is not allowed",
			slog.String("channel_id", binding.ChannelID.String()),
			slog.Int("status_code", response.StatusCode),
		)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("send YouTube live chat message: %w", youtubeAPIStatusError{statusCode: response.StatusCode})
	}

	return nil
}

func (c *ChatClient) cacheLiveChatID(ctx context.Context, key, liveChatID string, ttl time.Duration) {
	if err := c.liveChatIDCache.Set(ctx, key, liveChatID, kvoptions.WithExpire(ttl)); err != nil {
		c.logger.WarnContext(
			ctx,
			"youtube live chat cache write failed",
			slog.String("key", key),
			slog.Any("error", err),
		)
	}
}

func (c *ChatClient) invalidateLiveChatID(ctx context.Context, channelID string) {
	key := liveChatIDCachePrefix + channelID
	if err := c.liveChatIDCache.Delete(ctx, key); err != nil {
		c.logger.WarnContext(
			ctx,
			"youtube live chat cache invalidation failed",
			slog.String("key", key),
			slog.Any("error", err),
		)
	}
}

func (c *ChatClient) endpoint(path string) string {
	return strings.TrimRight(c.apiBaseURL, "/") + path
}

func splitMessage(text string) []string {
	normalizedText := strings.ReplaceAll(text, "\n", " ")
	if normalizedText == "" {
		return nil
	}

	runes := []rune(normalizedText)
	parts := make([]string, 0, (len(runes)/chatMessageMaxRunes)+1)
	for len(runes) > 0 {
		end := min(len(runes), chatMessageMaxRunes)
		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}

	return parts
}

func noActiveBroadcastError(binding channelplatformentity.ChannelPlatform) error {
	return fmt.Errorf("no active YouTube broadcast with live chat for channel %q", binding.PlatformChannelID)
}
