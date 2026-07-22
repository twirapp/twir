package twitch

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	buscore "github.com/twirapp/twir/libs/bus-core"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	twitchclient "github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TwirBus            *buscore.Bus
	Config             config.Config
	CachedTwitchClient *twitchcahe.CachedTwitchClient
	UsersRepository    users.Repository
	ChannelService     *channelservice.ChannelService
}

func New(opts Opts) *Service {
	return &Service{
		twirBus:            opts.TwirBus,
		config:             opts.Config,
		cachedTwitchClient: opts.CachedTwitchClient,
		usersRepository:    opts.UsersRepository,
		channelService:     opts.ChannelService,
	}
}

type Service struct {
	twirBus            *buscore.Bus
	config             config.Config
	cachedTwitchClient *twitchcahe.CachedTwitchClient
	usersRepository    users.Repository
	channelService     channelLookup
	newAppClient       twitchAppClientFactory
	newUserClient      twitchUserClientFactory
}

type channelLookup interface {
	GetChannelByID(ctx context.Context, id uuid.UUID) (channelsmodel.Channel, error)
}

type twitchAppClientFactory func(context.Context) (*helix.Client, error)

type twitchUserClientFactory func(context.Context, uuid.UUID) (*helix.Client, error)

func (c *Service) createAppClient(ctx context.Context) (*helix.Client, error) {
	if c.newAppClient != nil {
		return c.newAppClient(ctx)
	}

	return twitchclient.NewAppClientWithContext(ctx, c.config, c.twirBus)
}

func (c *Service) createUserClient(ctx context.Context, userID uuid.UUID) (*helix.Client, error) {
	if c.newUserClient != nil {
		return c.newUserClient(ctx, userID)
	}

	return twitchclient.NewUserClientWithContext(ctx, userID, c.config, c.twirBus)
}
