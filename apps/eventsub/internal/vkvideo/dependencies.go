package vkvideo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

const (
	deduplicationTTL    = 24 * time.Hour
	deduplicationPrefix = "eventsub:vk-video-live:message:"
)

type tokenProvider interface {
	GetUserToken(context.Context, uuid.UUID) (string, error)
}

type webSocketTokenProvider interface {
	DiscoverChatChannel(context.Context, uuid.UUID) (string, error)
	ConnectionToken(context.Context, uuid.UUID) (string, error)
	SubscriptionToken(context.Context, uuid.UUID, string) (string, error)
}

type userStore interface {
	GetByPlatformID(context.Context, platformentity.Platform, string) (usersmodel.User, error)
	Create(context.Context, usersrepository.CreateInput) (usersmodel.User, error)
}

type messagePublisher interface {
	Publish(context.Context, generic.ChatMessage) error
}

type messageDeduplicator interface {
	Claim(context.Context, string) (bool, error)
}

type busTokenProvider struct {
	request interface {
		Request(context.Context, buscoretokens.GetUserTokenRequest) (*buscore.QueueResponse[buscoretokens.TokenResponse], error)
	}
}

func (p busTokenProvider) GetUserToken(ctx context.Context, userID uuid.UUID) (string, error) {
	response, err := p.request.Request(ctx, buscoretokens.GetUserTokenRequest{UserId: userID})
	if err != nil {
		return "", fmt.Errorf("request VK Video user token: %w", err)
	}
	return response.Data.AccessToken, nil
}

type devAPIWebSocketTokenProvider struct {
	oauthTokens tokenProvider
	client      *vk.WebSocketTokenClient
}

func (p devAPIWebSocketTokenProvider) DiscoverChatChannel(ctx context.Context, userID uuid.UUID) (string, error) {
	oauthAccessToken, err := p.oauthTokens.GetUserToken(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get VK Video OAuth token: %w", err)
	}

	channel, err := p.client.DiscoverChatChannel(ctx, vk.OAuthAccessToken(oauthAccessToken))
	if err != nil {
		return "", fmt.Errorf("discover VK Video chat channel: %w", err)
	}

	return string(channel), nil
}

func (p devAPIWebSocketTokenProvider) ConnectionToken(ctx context.Context, userID uuid.UUID) (string, error) {
	oauthAccessToken, err := p.oauthTokens.GetUserToken(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get VK Video OAuth token: %w", err)
	}

	token, err := p.client.ConnectionToken(ctx, vk.OAuthAccessToken(oauthAccessToken))
	if err != nil {
		return "", fmt.Errorf("get VK Video WebSocket connection token: %w", err)
	}

	return string(token), nil
}

func (p devAPIWebSocketTokenProvider) SubscriptionToken(
	ctx context.Context,
	userID uuid.UUID,
	channel string,
) (string, error) {
	oauthAccessToken, err := p.oauthTokens.GetUserToken(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get VK Video OAuth token: %w", err)
	}

	token, err := p.client.SubscriptionToken(ctx, vk.OAuthAccessToken(oauthAccessToken), vk.WebSocketChannel(channel))
	if err != nil {
		return "", fmt.Errorf("get VK Video WebSocket subscription token: %w", err)
	}

	return string(token), nil
}

type redisDeduplicator struct {
	redis *redis.Client
}

func (d redisDeduplicator) Claim(ctx context.Context, id string) (bool, error) {
	claimed, err := d.redis.SetNX(ctx, deduplicationPrefix+id, "1", deduplicationTTL).Result()
	if err != nil {
		return false, fmt.Errorf("claim VK Video message id: %w", err)
	}
	return claimed, nil
}
