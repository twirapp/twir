package channel

import (
	"context"
	"log/slog"

	botplatforms "github.com/twirapp/twir/apps/bots/internal/platforms"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"gorm.io/gorm"
)

type twitchActionsClient interface {
	Ban(context.Context, twitchactions.BanOpts) error
	DeleteMessage(context.Context, twitchactions.DeleteMessageOpts) error
}

func New(
	logger *slog.Logger,
	gormDB *gorm.DB,
	twitchActions *twitchactions.TwitchActions,
	chatRegistry *platformsregistry.Registry[botplatforms.ChatAdapter],
	workersPool *workers.Pool,
	channelService *channelservice.ChannelService,
	usersRepo usersrepository.Repository,
) *Service {
	return &Service{
		gorm:           gormDB,
		logger:         logger,
		twitchActions:  twitchActions,
		chatRegistry:   chatRegistry,
		workersPool:    workersPool,
		channelService: channelService,
		usersRepo:      usersRepo,
	}
}

type Service struct {
	logger         *slog.Logger
	gorm           *gorm.DB
	twitchActions  twitchActionsClient
	chatRegistry   *platformsregistry.Registry[botplatforms.ChatAdapter]
	workersPool    *workers.Pool
	channelService *channelservice.ChannelService
	usersRepo      usersrepository.Repository
}
