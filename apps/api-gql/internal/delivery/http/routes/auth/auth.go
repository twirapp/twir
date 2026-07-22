package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/kv"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	botsrepo "github.com/twirapp/twir/libs/repositories/bots"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
	"github.com/twirapp/twir/libs/repositories/tokens"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Huma                 huma.API
	Config               config.Config
	Bus                  *buscore.Bus
	Sessions             *sessions.Auth
	Logger               *slog.Logger
	TrmManager           trm.Manager
	PlatformRegistry     *appplatform.Registry
	TokensRepository     tokens.Repository
	ChannelsRepo         channelsrepo.Repository
	ChannelPlatformsRepo channelplatformsrepo.Repository
	BotsRepo             botsrepo.Repository
	UsersRepo            usersrepository.Repository
	KickProvider         *kickplatform.Provider
	KickBotsRepo         kickbotsrepo.Repository
	KV                   kv.KV
}

type Auth struct {
	config                 config.Config
	bus                    *buscore.Bus
	sessions               sessionStore
	logger                 *slog.Logger
	transactionRunner      transactionRunner
	platformRegistry       *appplatform.Registry
	tokensRepository       tokens.Repository
	channelsRepo           channelsrepo.Repository
	channelPlatformsRepo   channelplatformsrepo.Repository
	botsRepo               botsrepo.Repository
	usersRepo              usersrepository.Repository
	kickProvider           *kickplatform.Provider
	kickBotsRepo           kickbotsrepo.Repository
	kv                     kv.KV
	eventSubPublisher      eventSubPublisher
	bindingConfigResolvers map[platformentity.Platform]platformBindingConfigResolver
	postPlatformAuthHooks  map[platformentity.Platform]postPlatformAuthHook
}

type sessionStore interface {
	GetInternalUserID(context.Context) (uuid.UUID, error)
	SetSessionInternalUserID(context.Context, uuid.UUID) error
	SetSessionCurrentPlatform(context.Context, string) error
	SetSessionSelectedDashboard(context.Context, string) error
	SetSessionTwitchUser(context.Context, helix.User) error
	SetSessionKickUser(context.Context, sessions.KickSessionUser) error
	SetOAuthAttempt(context.Context, string, sessions.OAuthAttempt) error
	GetOAuthAttempt(context.Context, string) (sessions.OAuthAttempt, error)
	DeleteOAuthAttempt(context.Context, string) error
}

func New(opts Opts) *Auth {
	p := &Auth{
		config:               opts.Config,
		bus:                  opts.Bus,
		sessions:             opts.Sessions,
		logger:               opts.Logger,
		transactionRunner:    opts.TrmManager,
		platformRegistry:     opts.PlatformRegistry,
		tokensRepository:     opts.TokensRepository,
		channelsRepo:         opts.ChannelsRepo,
		channelPlatformsRepo: opts.ChannelPlatformsRepo,
		botsRepo:             opts.BotsRepo,
		usersRepo:            opts.UsersRepo,
		kickProvider:         opts.KickProvider,
		kickBotsRepo:         opts.KickBotsRepo,
		kv:                   opts.KV,
	}
	if opts.Bus != nil && opts.Bus.EventSub != nil {
		p.eventSubPublisher = opts.Bus.EventSub.SubscribeToAllEvents
	}
	p.bindingConfigResolvers = map[platformentity.Platform]platformBindingConfigResolver{
		platformentity.PlatformTwitch: p.twitchBindingConfig,
		platformentity.PlatformKick:   p.kickBindingConfig,
	}
	p.postPlatformAuthHooks = map[platformentity.Platform]postPlatformAuthHook{
		platformentity.PlatformKick: p.updateKickBotTokenAfterAuth,
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
			OperationID: "auth-platform-code",
			Method:      http.MethodPost,
			Path:        "/auth/{platform}/code",
			Tags:        []string{"Auth"},
			Summary:     "Platform OAuth code exchange",
		},
		func(
			ctx context.Context, i *struct {
				Platform platformentity.Platform `path:"platform"`
				Body     platformCodeBody
			},
		) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			return p.handlePlatformCode(ctx, platformCodeInput{
				Platform: i.Platform,
				Code:     i.Body.Code,
				State:    i.Body.State,
				DeviceID: i.Body.DeviceID,
			})
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
			OperationID: "auth-platform-authorize",
			Method:      http.MethodGet,
			Path:        "/auth/{platform}/authorize",
			Tags:        []string{"Auth"},
			Summary:     "Get platform OAuth authorize URL",
		},
		func(ctx context.Context, i *struct {
			Platform   platformentity.Platform `path:"platform"`
			RedirectTo string                  `query:"redirect_to"`
		}) (*kickAuthorizeOutput, error) {
			return p.handlePlatformAuthorize(ctx, i.Platform, i.RedirectTo)
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
