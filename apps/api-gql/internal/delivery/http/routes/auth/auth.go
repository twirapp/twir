package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/kv"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	botsrepo "github.com/twirapp/twir/libs/repositories/bots"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
	"github.com/twirapp/twir/libs/repositories/tokens"
	userplatformaccountsrepo "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Huma                     huma.API
	Gorm                     *gorm.DB
	Config                   config.Config
	Bus                      *buscore.Bus
	Sessions                 *sessions.Auth
	Logger                   *slog.Logger
	TokensRepository         tokens.Repository
	UserPlatformAccountsRepo userplatformaccountsrepo.Repository
	ChannelsRepo             channelsrepo.Repository
	BotsRepo                 botsrepo.Repository
	UsersRepo                usersrepository.Repository
	KickProvider             *kickplatform.Provider
	KickBotsRepo             kickbotsrepo.Repository
	KV                       kv.KV
}

type Auth struct {
	gorm                     *gorm.DB
	config                   config.Config
	bus                      *buscore.Bus
	sessions                 *sessions.Auth
	logger                   *slog.Logger
	tokensRepository         tokens.Repository
	userPlatformAccountsRepo userplatformaccountsrepo.Repository
	channelsRepo             channelsrepo.Repository
	botsRepo                 botsrepo.Repository
	usersRepo                usersrepository.Repository
	kickProvider             *kickplatform.Provider
	kickBotsRepo             kickbotsrepo.Repository
	kv                       kv.KV
}

func New(opts Opts) *Auth {
	p := &Auth{
		gorm:                     opts.Gorm,
		config:                   opts.Config,
		bus:                      opts.Bus,
		sessions:                 opts.Sessions,
		logger:                   opts.Logger,
		tokensRepository:         opts.TokensRepository,
		userPlatformAccountsRepo: opts.UserPlatformAccountsRepo,
		channelsRepo:             opts.ChannelsRepo,
		botsRepo:                 opts.BotsRepo,
		usersRepo:                opts.UsersRepo,
		kickProvider:             opts.KickProvider,
		kickBotsRepo:             opts.KickBotsRepo,
		kv:                       opts.KV,
	}

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-post-code",
			Method:      http.MethodPost,
			Path:        "/auth",
			Tags:        []string{"Auth"},
			Summary:     "Auth post code",
		},
		func(
			ctx context.Context, i *struct {
				Body authBody
			},
		) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			return p.handleAuthPostCode(ctx, i.Body)
		},
	)

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-kick-code",
			Method:      http.MethodPost,
			Path:        "/auth/kick/code",
			Tags:        []string{"Auth"},
			Summary:     "Kick OAuth code exchange",
		},
		func(
			ctx context.Context, i *struct {
				Body kickCodeBody
			},
		) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			return p.handleKickCode(ctx, i.Body)
		},
	)

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-kick-authorize",
			Method:      http.MethodGet,
			Path:        "/auth/kick/authorize",
			Tags:        []string{"Auth"},
			Summary:     "Get Kick OAuth authorize URL",
		},
		func(ctx context.Context, i *struct {
			RedirectTo string `query:"redirect_to"`
		},
		) (*kickAuthorizeOutput, error) {
			return p.handleKickAuthorize(ctx, kickAuthorizeInput{RedirectTo: i.RedirectTo})
		},
	)

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "auth-kick-bot-callback",
			Method:      http.MethodGet,
			Path:        "/auth/kick/bot-callback",
			Tags:        []string{"Auth"},
			Summary:     "Kick bot setup callback",
		},
		func(ctx context.Context, i *struct {
			Code  string `query:"code"`
			State string `query:"state"`
		},
		) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			return p.handleKickBotCallback(ctx, kickBotCallbackInput{Code: i.Code, State: i.State})
		},
	)

	return p
}
