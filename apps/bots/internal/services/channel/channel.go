package channel

import (
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/kick"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger         *slog.Logger
	Gorm           *gorm.DB
	KickChatClient *kick.ChatClient
	TwitchActions  *twitchactions.TwitchActions
	WorkersPool    *workers.Pool
	ChannelService *channelservice.ChannelService
	UsersRepo      usersrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		gorm:           opts.Gorm,
		kickChatClient: opts.KickChatClient,
		logger:         opts.Logger,
		twitchActions:  opts.TwitchActions,
		workersPool:    opts.WorkersPool,
		channelService: opts.ChannelService,
		usersRepo:      opts.UsersRepo,
	}
}

type Service struct {
	logger         *slog.Logger
	gorm           *gorm.DB
	kickChatClient *kick.ChatClient
	twitchActions  *twitchactions.TwitchActions
	workersPool    *workers.Pool
	channelService *channelservice.ChannelService
	usersRepo      usersrepository.Repository
}
