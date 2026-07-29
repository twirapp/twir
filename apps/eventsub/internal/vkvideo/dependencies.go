package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	oauthvkvideo "github.com/twirapp/twir/libs/oauth/vkvideo"
	"github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
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

type chatUserEnsurer interface {
	UnsureUser(context.Context, user_creator.CreateUserInput) (*usersmodel.User, *usersstatsmodel.UserStat, error)
}

type messagePublisher interface {
	Publish(context.Context, generic.ChatMessage) error
}

type messageDeduplicator interface {
	Claim(context.Context, string) (bool, error)
}

type runtimeTokenProvider struct {
	source oauthvkvideo.UserTokenSource
}

func (p runtimeTokenProvider) GetUserToken(ctx context.Context, userID uuid.UUID) (string, error) {
	credential, err := p.source.Token(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("load VK Video user token: %w", err)
	}
	return credential.AccessToken, nil
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

type bindingsProvider func(context.Context) ([]channelplatformentity.ChannelPlatform, error)

func newDatabaseBindingsProvider(repo channels.Repository) bindingsProvider {
	return func(ctx context.Context) ([]channelplatformentity.ChannelPlatform, error) {
		if repo == nil {
			return nil, errors.New("channels repository is not configured")
		}

		channelList, err := repo.GetAllByBindingPlatform(ctx, platformentity.PlatformVKVideoLive)
		if err != nil {
			return nil, fmt.Errorf("list VK Video Live channels: %w", err)
		}

		bindings := make([]channelplatformentity.ChannelPlatform, 0, len(channelList))
		for _, channel := range channelList {
			binding, ok := channel.Binding(platformentity.PlatformVKVideoLive)
			if !ok || !binding.Enabled {
				continue
			}
			bindings = append(bindings, binding)
		}

		return bindings, nil
	}
}
