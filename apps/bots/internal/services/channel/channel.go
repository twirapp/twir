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
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger         *slog.Logger
	Gorm           *gorm.DB
	TwitchActions  *twitchactions.TwitchActions
	ChatRegistry   *platformsregistry.Registry[botplatforms.ChatAdapter]
	WorkersPool    *workers.Pool
	ChannelService *channelservice.ChannelService
	UsersRepo      usersrepository.Repository
}

type twitchActionsClient interface {
	Ban(context.Context, twitchactions.BanOpts) error
	DeleteMessage(context.Context, twitchactions.DeleteMessageOpts) error
}

func New(opts Opts) *Service {
	return &Service{
		gorm:           opts.Gorm,
		logger:         opts.Logger,
		twitchActions:  opts.TwitchActions,
		chatRegistry:   opts.ChatRegistry,
		workersPool:    opts.WorkersPool,
		channelService: opts.ChannelService,
		usersRepo:      opts.UsersRepo,
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
