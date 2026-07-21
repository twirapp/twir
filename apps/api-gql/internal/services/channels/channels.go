package channels

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	chatmessagesrepo "github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsRepository channels.Repository
	UsersRepository    users.Repository
	ChannelService     *channelservice.ChannelService
}

func New(opts Opts) *Service {
	return &Service{
		channelsRepository: opts.ChannelsRepository,
		usersRepository:    opts.UsersRepository,
		channelService:     opts.ChannelService,
	}
}

type Service struct {
	channelsRepository channels.Repository
	usersRepository    users.Repository
	channelService     *channelservice.ChannelService
}

var ErrNotFound = fmt.Errorf("channel not found")

func (c *Service) mapToEntity(m model.Channel) entity.Channel {
	return entity.Channel{
		ID:             m.ID,
		IsEnabled:      m.IsEnabled,
		IsTwitchBanned: m.IsTwitchBanned,
		IsBotMod:       m.IsBotMod,
		BotID:          m.BotID,
	}
}

func (c *Service) GetByID(ctx context.Context, channelID uuid.UUID) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByID(ctx, channelID)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel), nil
}

func (c *Service) GetByTwitchPlatformID(ctx context.Context, twitchPlatformID string) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByPlatformUserID(ctx, twitchPlatformID, platformentity.PlatformTwitch)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel), nil
}

func (c *Service) GetByKickPlatformID(ctx context.Context, kickPlatformID string) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByPlatformUserID(ctx, kickPlatformID, platformentity.PlatformKick)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel), nil
}

type ApiKeyChannelIdentity struct {
	InternalChannelID string
	ChatTargets       []chatmessagesrepo.PlatformChannelIdentity
}

type ChannelPlatformIdentity struct {
	Platform platformentity.Platform
	ID       string
}

func (c *Service) ResolveApiKeyChannelIdentityByAnyPlatformUUID(ctx context.Context, userId uuid.UUID) (*ApiKeyChannelIdentity, error) {
	user, err := c.usersRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	var channel model.Channel

	switch user.Platform {
	case platformentity.PlatformKick:
		channel, err = c.channelService.GetChannelByConnectedUser(ctx, user.ID, platformentity.PlatformKick)
		if err != nil {
			return nil, fmt.Errorf("failed to get kick channel: %w", err)
		}
	default:
		channel, err = c.channelService.GetChannelByConnectedUser(ctx, user.ID, platformentity.PlatformTwitch)
		if err != nil {
			return nil, fmt.Errorf("failed to get twitch channel: %w", err)
		}
	}

	targets := make([]chatmessagesrepo.PlatformChannelIdentity, 0, 2)
	for _, identity := range c.mapChannelPlatformIdentities(channel) {
		targets = append(
			targets, chatmessagesrepo.PlatformChannelIdentity{
				Platform:          identity.Platform.String(),
				PlatformChannelID: identity.ID,
			},
		)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no chat message targets found for api key")
	}

	return &ApiKeyChannelIdentity{
		InternalChannelID: channel.ID.String(),
		ChatTargets:       targets,
	}, nil
}

func (c *Service) ResolveApiKeyChannelIdentityByUserOrChannelApiKey(
	ctx context.Context,
	apiKey string,
) (ApiKeyChannelIdentity, error) {
	var channel model.Channel
	foundedChannel, err := c.channelsRepository.GetByApiKey(ctx, apiKey)
	if err != nil && !errors.Is(err, channels.ErrNotFound) {
		return ApiKeyChannelIdentity{}, err
	}

	if !foundedChannel.IsNil() {
		channel = foundedChannel
	} else {
		user, err := c.usersRepository.GetByApiKey(ctx, apiKey)
		if err != nil {
			return ApiKeyChannelIdentity{}, fmt.Errorf("failed to get user: %w", err)
		}

		switch user.Platform {
		case platformentity.PlatformKick:
			channel, err = c.channelService.GetChannelByConnectedUser(ctx, user.ID, platformentity.PlatformKick)
			if err != nil {
				return ApiKeyChannelIdentity{}, fmt.Errorf("failed to get kick channel: %w", err)
			}
		default:
			channel, err = c.channelService.GetChannelByConnectedUser(ctx, user.ID, platformentity.PlatformTwitch)
			if err != nil {
				return ApiKeyChannelIdentity{}, fmt.Errorf("failed to get twitch channel: %w", err)
			}
		}
	}

	targets := make([]chatmessagesrepo.PlatformChannelIdentity, 0, 2)
	for _, identity := range c.mapChannelPlatformIdentities(channel) {
		targets = append(
			targets, chatmessagesrepo.PlatformChannelIdentity{
				Platform:          identity.Platform.String(),
				PlatformChannelID: identity.ID,
			},
		)
	}

	if len(targets) == 0 {
		return ApiKeyChannelIdentity{}, fmt.Errorf("no chat message targets found for api key")
	}

	return ApiKeyChannelIdentity{
		InternalChannelID: channel.ID.String(),
		ChatTargets:       targets,
	}, nil
}

func (c *Service) GetPlatformIdentities(ctx context.Context, channelID uuid.UUID) ([]ChannelPlatformIdentity, error) {
	channel, err := c.channelService.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	return c.mapChannelPlatformIdentities(channel), nil
}

func (c *Service) mapChannelPlatformIdentities(channel model.Channel) []ChannelPlatformIdentity {
	identities := make([]ChannelPlatformIdentity, 0, 2)
	if channel.TwitchPlatformID != nil && *channel.TwitchPlatformID != "" {
		identities = append(
			identities, ChannelPlatformIdentity{
				Platform: platformentity.PlatformTwitch,
				ID:       *channel.TwitchPlatformID,
			},
		)
	}

	if channel.KickPlatformID != nil && *channel.KickPlatformID != "" {
		identities = append(
			identities, ChannelPlatformIdentity{
				Platform: platformentity.PlatformKick,
				ID:       *channel.KickPlatformID,
			},
		)
	}

	return identities
}
