package twitch

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	twitchclient "github.com/twirapp/twir/libs/twitch"
)

func New(
	twirBus *buscore.Bus,
	cfg config.Config,
	cachedTwitchClient *twitchcahe.CachedTwitchClient,
	usersRepository users.Repository,
	channelService *channelservice.ChannelService,
	kickProvider *kickplatform.Provider,
) *Service {
	return &Service{
		twirBus:            twirBus,
		config:             cfg,
		cachedTwitchClient: cachedTwitchClient,
		usersRepository:    usersRepository,
		channelService:     channelService,
		kickProvider:       kickProvider,
	}
}

type Service struct {
	twirBus                     *buscore.Bus
	config                      config.Config
	cachedTwitchClient          *twitchcahe.CachedTwitchClient
	usersRepository             users.Repository
	channelService              channelLookup
	kickProvider                *kickplatform.Provider
	newAppClient                twitchAppClientFactory
	newUserClient               twitchUserClientFactory
	requestKickUserToken        kickUserTokenRequester
	updateKickStreamInformation kickStreamInformationUpdater
}

type channelLookup interface {
	GetChannelByID(ctx context.Context, id uuid.UUID) (channelentity.Channel, error)
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
