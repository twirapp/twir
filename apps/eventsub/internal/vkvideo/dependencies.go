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
