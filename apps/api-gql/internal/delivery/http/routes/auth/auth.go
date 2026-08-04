package auth

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/kv"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	kickplatform "github.com/twirapp/twir/apps/api-gql/internal/platform/kick"
	vkvideo "github.com/twirapp/twir/apps/api-gql/internal/platform/vkvideo"
	youtube "github.com/twirapp/twir/apps/api-gql/internal/platform/youtube"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	botsrepo "github.com/twirapp/twir/libs/repositories/bots"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
	"github.com/twirapp/twir/libs/repositories/tokens"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	vkvideobotsrepo "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	youtubebotsrepo "github.com/twirapp/twir/libs/repositories/youtube_bots"
)

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
	vkVideoBotProvider     vkVideoBotSetupProvider
	vkVideoBotsRepo        vkvideobotsrepo.Repository
	youtubeBotProvider     youtubeBotSetupProvider
	youtubeBotsRepo        youtubebotsrepo.Repository
	kv                     kv.KV
	vkVideoBotSetupStateMu sync.Mutex
	youtubeBotSetupStateMu sync.Mutex
	dashboardAccess        dashboardAccessChecker
	eventSubPublisher      eventSubPublisher
	bindingConfigResolvers map[platformentity.Platform]platformBindingConfigResolver
	postPlatformAuthHooks  map[platformentity.Platform]postPlatformAuthHook
}

type youtubeBotCallbackOutput struct {
	Status   int
	Location string `header:"Location"`
}

type dashboardAccessChecker interface {
	IsOwner(context.Context, string, uuid.UUID) (bool, error)
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

func New(humaAPI huma.API, config config.Config, bus *buscore.Bus, sessionAuth *sessions.Auth, logger *slog.Logger, trmManager trm.Manager, platformRegistry *appplatform.Registry, tokensRepository tokens.Repository, channelsRepo channelsrepo.Repository, channelPlatformsRepo channelplatformsrepo.Repository, botsRepo botsrepo.Repository, usersRepo usersrepository.Repository, kickProvider *kickplatform.Provider, kickBotsRepo kickbotsrepo.Repository, vkVideoBotProvider *vkvideo.BotSetupProvider, vkVideoBotsRepo vkvideobotsrepo.Repository, youtubeBotProvider *youtube.Provider, youtubeBotsRepo youtubebotsrepo.Repository, kvClient kv.KV, dashboardAccess *dashboardaccess.Service) *Auth {
	p := &Auth{
		config:               config,
		bus:                  bus,
		sessions:             sessionAuth,
		logger:               logger,
		transactionRunner:    trmManager,
		platformRegistry:     platformRegistry,
		tokensRepository:     tokensRepository,
		channelsRepo:         channelsRepo,
		channelPlatformsRepo: channelPlatformsRepo,
		botsRepo:             botsRepo,
		usersRepo:            usersRepo,
		kickProvider:         kickProvider,
		kickBotsRepo:         kickBotsRepo,
		vkVideoBotProvider:   vkVideoBotProvider,
		vkVideoBotsRepo:      vkVideoBotsRepo,
		youtubeBotProvider:   youtubeBotProvider,
		youtubeBotsRepo:      youtubeBotsRepo,
		kv:                   kvClient,
		dashboardAccess:      dashboardAccess,
	}
	if bus != nil && bus.EventSub != nil {
		p.eventSubPublisher = bus.EventSub.SubscribeToAllEvents
	}
	p.bindingConfigResolvers = map[platformentity.Platform]platformBindingConfigResolver{
		platformentity.PlatformTwitch:      p.twitchBindingConfig,
		platformentity.PlatformKick:        p.kickBindingConfig,
		platformentity.PlatformVKVideoLive: p.vkVideoBotBindingConfig,
		platformentity.PlatformYouTube:     p.youtubeBotBindingConfig,
	}
	p.postPlatformAuthHooks = map[platformentity.Platform]postPlatformAuthHook{
		platformentity.PlatformKick: p.updateKickBotTokenAfterAuth,
	}

	huma.Register(
		humaAPI,
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
		humaAPI,
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
		humaAPI,
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
		humaAPI,
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
		humaAPI,
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
		humaAPI,
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

	huma.Register(
		humaAPI,
		huma.Operation{
			OperationID: "auth-vk-video-bot-callback",
			Method:      http.MethodGet,
			Path:        "/auth/vk-video/bot-callback",
			Tags:        []string{"Auth"},
			Summary:     "VK Video bot setup callback",
		},
		func(ctx context.Context, i *struct {
			Code  string `query:"code"`
			State string `query:"state"`
		}) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
			if err := p.CompleteVKVideoBotSetup(ctx, i.Code, i.State); err != nil {
				return nil, huma.Error400BadRequest("Cannot complete VK Video bot setup", err)
			}

			return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: p.config.GetVkVideoBotCallbackUrl()}), nil
		},
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			OperationID: "auth-youtube-bot-callback",
			Method:      http.MethodGet,
			Path:        "/auth/youtube/bot-callback",
			Tags:        []string{"Auth"},
			Summary:     "YouTube bot setup callback",
		},
		func(ctx context.Context, i *struct {
			Code  string `query:"code"`
			State string `query:"state"`
		}) (*youtubeBotCallbackOutput, error) {
			return p.completeYouTubeBotCallback(ctx, i.Code, i.State)
		},
	)

	return p
}

func (a *Auth) completeYouTubeBotCallback(ctx context.Context, code, state string) (*youtubeBotCallbackOutput, error) {
	if err := a.CompleteYouTubeBotSetup(ctx, code, state); err != nil {
		return nil, huma.Error400BadRequest("Cannot complete YouTube bot setup", err)
	}

	return &youtubeBotCallbackOutput{Status: http.StatusFound, Location: "/dashboard/admin"}, nil
}
