package channels

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels"
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

func (c *Service) mapToEntity(m channelentity.Channel, p platformentity.Platform) (entity.Channel, error) {
	result := entity.Channel{ID: m.ID}

	binding, found := m.Binding(p)
	if !found {
		return result, nil
	}

	result.IsEnabled = binding.Enabled
	botConfig, err := binding.ParseTwitchBotConfig()
	if err != nil {
		return entity.ChannelNil, fmt.Errorf("unmarshal %s channel bot config: %w", p, err)
	}

	result.BotID = botConfig.BotID
	result.IsBotMod = botConfig.IsBotMod
	result.IsTwitchBanned = botConfig.IsTwitchBanned
	return result, nil
}

func (c *Service) GetByID(ctx context.Context, channelID uuid.UUID) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByID(ctx, channelID)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel, platformentity.PlatformTwitch)
}

func (c *Service) GetByTwitchPlatformID(ctx context.Context, twitchPlatformID string) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platformentity.PlatformTwitch,
		twitchPlatformID,
	)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel, platformentity.PlatformTwitch)
}

func (c *Service) GetByKickPlatformID(ctx context.Context, kickPlatformID string) (entity.Channel, error) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platformentity.PlatformKick,
		kickPlatformID,
	)
	if err != nil {
		if err == channels.ErrNotFound {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}

	return c.mapToEntity(channel, platformentity.PlatformKick)
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

	var channel channelentity.Channel

	channel, err = c.channelService.GetChannelByBindingUserID(ctx, user.Platform, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s channel: %w", user.Platform, err)
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
	var channel channelentity.Channel
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

		channel, err = c.channelService.GetChannelByBindingUserID(ctx, user.Platform, user.ID)
		if err != nil {
			return ApiKeyChannelIdentity{}, fmt.Errorf("failed to get %s channel: %w", user.Platform, err)
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

func (c *Service) mapChannelPlatformIdentities(channel channelentity.Channel) []ChannelPlatformIdentity {
	identities := make([]ChannelPlatformIdentity, 0, len(channel.Bindings))
	for _, binding := range channel.Bindings {
		if binding.PlatformChannelID == "" {
			continue
		}

		identities = append(
			identities, ChannelPlatformIdentity{
				Platform: binding.Platform,
				ID:       binding.PlatformChannelID,
			},
		)
	}

	return identities
}
