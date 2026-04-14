package channel

import (
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/kick"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	userplatformaccountsrepository "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger                   *slog.Logger
	Gorm                     *gorm.DB
	KickChatClient           *kick.ChatClient
	TwitchActions            *twitchactions.TwitchActions
	WorkersPool              *workers.Pool
	ChannelsRepo             channelsrepository.Repository
	UserPlatformAccountsRepo userplatformaccountsrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		gorm:                     opts.Gorm,
		kickChatClient:           opts.KickChatClient,
		logger:                   opts.Logger,
		twitchActions:            opts.TwitchActions,
		workersPool:              opts.WorkersPool,
		channelsRepo:             opts.ChannelsRepo,
		userPlatformAccountsRepo: opts.UserPlatformAccountsRepo,
	}
}

type Service struct {
	logger                   *slog.Logger
	gorm                     *gorm.DB
	kickChatClient           *kick.ChatClient
	twitchActions            *twitchactions.TwitchActions
	workersPool              *workers.Pool
	channelsRepo             channelsrepository.Repository
	userPlatformAccountsRepo userplatformaccountsrepository.Repository
}
